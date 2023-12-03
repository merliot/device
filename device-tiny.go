//go:build tinygo

package device

type deviceOS struct {
}

func (d *Device) deviceOSInit() {
}

func (d *Device) Serve(thinger dean.Thinger) {
	server := dean.NewServer(thinger)
	server.Dial()
	server.Run()
}
