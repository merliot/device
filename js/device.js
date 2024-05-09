class DeviceBase {

	constructor(div, view) {
		this.div = div
		this.view = view // 0=full, 1=tile
		this.state = null;
	}

	setSender(sender) {
		this.sender = sender
	}

	visible() {
		const rect = this.div.getBoundingClientRect()
		const viewportLeft = window.pageXOffset
		const viewportRight = viewportLeft + window.innerWidth
		const viewportTop = window.pageYOffset
		const viewportBottom = viewportTop + window.innerHeight

		// Check if any part of the element is within the viewport
		return (
			rect.right >= viewportLeft &&
			rect.left <= viewportRight &&
			rect.top <= viewportBottom &&
			rect.bottom >= viewportTop
		)
	}

	send(msg) {
		this.sender.send(this.state.Id, msg)
	}

	open() {
		this.state.Online ? this.online() : this.offline()
	}

	close() {
		this.offline()
	}

	online() {
		document.body.classList.replace("offline", "online")
	}

	offline() {
		document.body.classList.replace("online", "offline")
	}

	handle(msg) {
		// drop msg
	}

	handleMsg(msg) {
		switch(msg.Path) {
		case "state":
			this.state = msg
			this.open()
			break
		case "online":
			this.state.Online = true
			this.online()
			break
		case "offline":
			this.state.Online = false
			this.offline()
			break
		default:
			this.handle(msg)
			break
		}
	}
}

class Conn {

	constructor(wsUrl) {
		this.wsUrl = wsUrl
		this.ws = null
		this.pingID = null
		this.timeoutID = null
		this.devices = {}
		this.setupVisibilityChange()
		this.connect()
	}

	ping() {
		if (this.ws.readyState === 1) {
			this.ws.send(JSON.stringify({Path: "ping"}))
		}
	}

	connect() {

		const url = new URL(this.wsUrl)
		const params = new URLSearchParams(url.search)
		const pingPeriod = params.get("ping-period") * 1000

		if (this.ws != null) {
			// webSocket is still open...wait for it to close
			return
		}

		console.log('ws connecting...')
		this.ws = new WebSocket(this.wsUrl)

		this.ws.onopen = () => {
			console.log('ws opened')
			// Start ping
			this.pingID = setInterval(() => this.ping(), pingPeriod)
			this.opened()
		};

		this.ws.onclose = () => {
			console.log('ws close');
			clearInterval(this.pingID);
			if (document.visibilityState === 'visible') {
				// Schedule reconnect in 2 seconds
				this.timeoutID = setTimeout(() => this.connect(), 2000);
			}
			this.ws = null
			this.close()
		};

		this.ws.onmessage = (event) => {
			var msg = JSON.parse(event.data)
			switch(msg.Path) {
				case "pong":
					return
				default:
					console.log(msg)
					this.mux(msg)
					break
			}
		}

	}

	disconnect() {
		clearInterval(this.pingID)
		clearTimeout(this.timeoutID)
		if (this.ws && this.ws.readyState === 1) {
			this.ws.close()
		}
	}

	setupVisibilityChange() {
		document.addEventListener('visibilitychange', () => {
			if (document.visibilityState === 'visible') {
				this.connect()
			} else {
				this.disconnect()
			}
		})
	}
}

class Single extends Conn {

	constructor(wsUrl, device) {
		super(wsUrl)
		this.device = device
		this.device.setSender(this)
	}

	opened() {
		this.ws.send(JSON.stringify({Tags: "", Path: "get/state"}))
	}

	close() {
		this.device.close()
	}

	mux(msg) {
		device.handleMsg(msg)
	}

	send(id, msg) {
		msg.Tags = ""
		this.ws.send(JSON.stringify(msg))
	}
}

class Trunk extends Conn {

	constructor(wsUrl) {
		super(wsUrl)
		this.devices = {}
		window.addEventListener('resize', () => { this.checkVisibility() })
	}

	checkVisibility() {
	}

	registerDevice(id, deviceInstance) {
		if (!(deviceInstance instanceof DeviceBase)) {
			console.error("deviceInstance is not an instance of DeviceBase")
			return
		}
		this.devices[id] = deviceInstance
	}

	getState(id) {
		if (this.ws.readyState === 1) {
			this.ws.send(JSON.stringify({Tags: id, Path: "get/state"}))
		}
	}

	opened() {
		for (let id in this.devices) {
			let device = this.devices[id]
			if (device.visible()) {
				this.getState(id)
			}
		}
	}

	close() {
		for (let id in this.devices) {
			let device = this.devices[id]
			device.close()
		}
	}

	mux(msg) {
		let device = this.devices[msg.Id]
		if (device !== undefined && device.visible()) {
			device.handleMsg(msg)
		}
	}

	send(id, msg) {
		msg.Tags = id
		this.ws.send(JSON.stringify(msg))
	}

}
