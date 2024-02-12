//go:build tinygo

package uart

import "machine"

type Uart struct {
	*machine.UART
}

func New() Uart {
	u := Uart{machine.UART0}
	u.Configure(machine.UARTConfig{
		TX:       machine.UART0_TX_PIN,
		RX:       machine.UART0_RX_PIN,
		BaudRate: 9600,
	})
	return u
}
