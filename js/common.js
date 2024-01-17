class WebSocketController {

	constructor(prefix, url) {
		this.prefix = prefix
		this.url = url
		this.state = null;
		this.webSocket = null;
		this.pingID = null;
		this.pingAlive = false;
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

		console.log(this.prefix, 'connecting...');
		if (this.webSocket != null) {
			console.log(this.prefix, "webSocket is not NULL!", this.webSocket.readyState)
		}
		this.webSocket = new WebSocket(this.url);

		this.webSocket.onopen = () => {
			console.log(this.prefix, 'open');
			if (this.webSocket.readyState !== 1) {
				console.log(this.prefix, "webSocket not READY!", this.webSocket.readyState)
			}
			this.webSocket.send(JSON.stringify({Path: "get/state"}));
			this.pingID = setInterval(() => this.ping(), pingPeriod)
		};

		this.webSocket.onclose = () => {
			console.log(this.prefix, 'close');
			this.close();
			clearInterval(this.pingID);
			if (document.visibilityState === 'visible') {
				this.timeoutID = setTimeout(() => this.initWebSocket(), 2000);
			}
			this.webSocket = null
		};

		this.webSocket.onmessage = (event) => {

			if (event.data == "pong") {
				//console.log(this.prefix, "PONG", new Date())
				this.pingAlive = true
				return
			}

			var msg = JSON.parse(event.data)
			console.log(this.prefix, msg)

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
		};

	}

	closeWebSocket() {
		clearInterval(this.pingID);
		clearTimeout(this.timeoutID)
		if (this.webSocket && this.webSocket.readyState === 1) {
			this.webSocket.close(3000, "leaving page");
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
		//if (!this.pingAlive) {
		//	console.log(this.prefix, "NOT ALIVE", new Date() - this.pingSent)
			// This waits for an ACK from server, but the server
			// may be gone, it may take a bit to close the websocket
			//this.webSocket.close()
			//clearInterval(this.pingID)
			//return
		//}
		if (this.webSocket.readyState === 1) {
			this.pingAlive = false
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
		if (this.stat !== null) {
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
