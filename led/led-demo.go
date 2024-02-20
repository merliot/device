//go:build !rpi && !tinygo

package led

type Led struct {
	State bool
}

func (l Led) Configure() {}
func (l Led) Set(bool)   {}
func (l Led) On()        {}
func (l Led) Off()       {}
