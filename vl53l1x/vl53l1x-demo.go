//go:build !rpi && !tinygo

package vl53l1x

type Vl53l1x struct {
}

func (v Vl53l1x) Configure() {
}

func (v Vl53l1x) Distance() (int32, bool) {
	return 0, true
}
