package main

import "github.com/merliot/device/target"

//go:generate go run main.go
func main() {
	target.GenTargetJS("../../js/deployTargetGpios.js")
}
