class WebSocketController {

	constructor(prefix, url) {
		this.prefix = prefix
		this.url = url
		this.state = null;
		this.webSocket = null;
		this.pingID = null;
		this.pingSent = null;
		this.timeoutID = null;
		this.stat = document.getElementById("status");
		this.setupVisibilityChange();
		this.initWebSocket();
	}

	initWebSocket() {

		const url = new URL(this.url);
		const params = new URLSearchParams(url.search);
		const pingPeriod = params.get("ping-period") * 1000;

		if (this.webSocket != null) {
			// webSocket is still open...wait for it to close
			return
		}

		console.log(this.prefix, 'connecting...');
		document.body.style.cursor = 'wait'
		this.webSocket = new WebSocket(this.url);

		this.webSocket.onopen = () => {
			console.log(this.prefix, 'open');
			this.webSocket.send(JSON.stringify({Path: "get/state"}));
			// Start ping
			this.pingID = setInterval(() => this.ping(), pingPeriod)
		};

		this.webSocket.onclose = () => {
			console.log(this.prefix, 'close');
			clearInterval(this.pingID);
			if (document.visibilityState === 'visible') {
				// Schedule reconnect in 2 seconds
				this.timeoutID = setTimeout(() => this.initWebSocket(), 2000);
			}
			this.webSocket = null
			this.close();
		};

		this.webSocket.onmessage = (event) => {

			if (event.data == "pong") {
				return
			}

			var msg = JSON.parse(event.data)
			console.log(this.prefix, msg)

			switch(msg.Path) {
				case "state":
					this.state = msg
					this.open()
					document.body.style.cursor = 'default'
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
		};

	}

	closeWebSocket() {
		clearInterval(this.pingID)
		clearTimeout(this.timeoutID)
		if (this.webSocket && this.webSocket.readyState === 1) {
			this.webSocket.close()
		}
	}

	setupVisibilityChange() {
		document.addEventListener('visibilitychange', () => {
			if (document.visibilityState === 'visible') {
				this.initWebSocket();
			} else {
				this.closeWebSocket();
			}
		});
	}

	ping() {
		if (this.webSocket.readyState === 1) {
			this.webSocket.send("ping")
			this.pingSent = new Date()
		}
	}

	open() {
		this.state.Online ? this.online() : this.offline()
	}

	close() {
		this.offline()
	}

	online() {
		if (this.stat !== null) {
			this.stat.innerHTML = ""
			this.stat.style.border = "none"
			this.stat.style.color = "none"
		}
	}

	offline() {
		if (this.stat !== null && document.visibilityState === 'visible') {
			this.stat.innerHTML = "Offline"
			this.stat.style.border = "solid"
			this.stat.style.color = "red"
		}
	}

	handle(msg) {
		// drop msg
	}
}

export { WebSocketController };
