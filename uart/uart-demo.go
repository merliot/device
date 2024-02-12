//go:build !rpi && !tinygo

package uart

type Uart struct {
}

func New() Uart {
	return Uart{}
}
