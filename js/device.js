class DeviceBase {
	constructor(div) {
		this.div = div
		this.state = null;
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

class Trunk {

	constructor(wsUrl) {
		this.wsUrl = wsUrl
		this.ws = null
		this.pingID = null
		this.timeoutID = null
		this.devices = {}
		this.setupVisibilityChange()
		this.connect()
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

	getState(tag) {
		if (this.ws.readyState === 1) {
			this.ws.send(JSON.stringify({Tags: tag, Path: "get/state"}))
		}
	}

	getStates() {
		for (let tag in this.devices) {
			let device = this.devices[tag]
			if (device.visible()) {
				this.getState(tag)
			}
		}
	}

	closeAll() {
		for (let tag in this.devices) {
			let device = this.devices[tag]
			device.close()
		}
	}

	popTag(msg) {
		const tags = msg.Tags.split('.')
		const tag = tags.shift()
		msg.Tags = tags.join('.')
		return tag
	}

	mux(msg) {
		let device = this.devices[msg.Tag]
		if (device !== undefined && device.visible()) {
			this.popTag(msg)
			device.handleMsg(msg)
		}
	}

	ping() {
		if (this.ws.readyState === 1) {
			this.ws.send(JSON.stringify({Tags: "", Path: "ping"}))
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
			this.getStates()
		};

		this.ws.onclose = () => {
			console.log('ws close');
			clearInterval(this.pingID);
			if (document.visibilityState === 'visible') {
				// Schedule reconnect in 2 seconds
				this.timeoutID = setTimeout(() => this.connect(), 2000);
			}
			this.ws = null
			this.closeAll()
		};

		this.ws.onmessage = (event) => {
			var msg = JSON.parse(event.data)
			console.log(msg)
			switch(msg.Path) {
				case "pong":
					return
				default:
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
		window.addEventListener('resize', () => { this.checkVisibility() })
		document.addEventListener('visibilitychange', () => {
			if (document.visibilityState === 'visible') {
				this.connect()
			} else {
				this.disconnect()
			}
		})
	}
}
