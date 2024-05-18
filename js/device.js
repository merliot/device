class DeviceBase {

	constructor(container, view, assets) {
		this.container = container
		this.view = view // 0=full, 1=tile
		this.assets = assets
		this.state = null;
	}

	setSender(sender) {
		this.sender = sender
	}

	visible() {
		const rect = this.container.getBoundingClientRect()
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

	send(path, msg) {
		const tags = this.state.Id
		this.sender.send(tags, path, msg)
	}

	open() {
		this.state.Online ? this.online() : this.offline()
	}

	close() {
		this.offline()
	}

	online() {
		this.container.classList.replace("offline", "online")
	}

	offline() {
		this.container.classList.replace("online", "offline")
	}

	handle(path, msg) {
		// drop msg
	}

	handleMsg(path, msg) {

		switch(path) {
		case "state":
			this.state = msg
			break
		}

		let prefix = "[" + this.state.Model + " " + this.state.Name + "]"
		console.log(prefix, path, msg)

		switch(path) {
		case "state":
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
			this.handle(path, msg)
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
		this.send("", "ping", {})
	}

	send(tags, path, msg) {
		const payload = btoa(JSON.stringify(msg))
		const pkt = JSON.stringify({Tags: tags, Path: path, Payload: payload})
		if (this.ws.readyState === 1) {
			this.ws.send(pkt)
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
			const pkt = JSON.parse(event.data)
			switch(pkt.Path) {
				case "pong":
					return
				default:
					this.mux(pkt)
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
		this.send("", "get/state", {})
	}

	close() {
		this.device.close()
	}

	send(tags, path, msg) {
		const payload = btoa(JSON.stringify(msg))
		const pkt = JSON.stringify({Tags: "", Path: path, Payload: payload})
		if (this.ws.readyState === 1) {
			this.ws.send(pkt)
		}
	}

	mux(pkt) {
		const msg = JSON.parse(atob(pkt.Payload))
		this.device.handleMsg(pkt.Path, msg)
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

	registerDevice(tag, deviceInstance) {
		if (!(deviceInstance instanceof DeviceBase)) {
			console.error("deviceInstance is not an instance of DeviceBase")
			return
		}
		this.devices[tag] = deviceInstance
	}

	opened() {
		for (let tag in this.devices) {
			let device = this.devices[tag]
			if (device.visible()) {
				this.send(tag, "get/state", {})
			}
		}
	}

	close() {
		for (let tag in this.devices) {
			let device = this.devices[tag]
			device.close()
		}
	}

	popTag(msg) {
		const tags = msg.Tags.split(".")
		const tag = tags.shift()
		msg.Tags = tags.join(".")
		return tag
	}

	mux(pkt) {
		const tag = this.popTag(pkt)
		const device = this.devices[tag]
		if (device !== undefined && device.visible()) {
			const msg = JSON.parse(atob(pkt.Payload))
			device.handleMsg(pkt.Path, msg)
		}
	}
}
