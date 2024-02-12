//go:build tinygo

package uart

import "machine"

type Uart struct {
	*machine.UART
}

func New() Uart {
	u := Uart{machine.DefaultUART}
	u.Configure(machine.UARTConfig{
		TX:       machine.UART_TX_PIN,
		RX:       machine.UART_RX_PIN,
		BaudRate: 9600,
	})
	return u
}
