class DeviceGadget extends DeviceBase {

	open() {
		this.bottles = document.getElementById("bottles")
		this.bottles.textContent = this.state.Bottles
		document.getElementById("takeone").onclick = () => {
			this.send("takeone", {})
		}
	}

	handle(path, msg) {
		switch(path) {
		case "tookone":
			this.bottles.textContent = msg.Bottles
			break
		}
	}
}
