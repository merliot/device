package modbus

import "io"

type Transport struct {
	io.ReadWriter
}

func NewTransport(rw io.ReadWriter) Transport {
	return Transport{ReadWriter: rw}
}
