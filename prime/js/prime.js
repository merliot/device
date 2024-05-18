class DevicePrime extends DeviceBase {

	open() {
		super.open()
		document.title = this.state.Child.Model + " - " + this.state.Child.Name
	}
}
