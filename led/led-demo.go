//go:build !rpi && !tinygo

package led

type Led struct {
	State bool
}

func New() Led {
	return Led{}
}

func (l Led) Set(state bool) {
}
