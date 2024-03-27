// go run ./cmd
// go run -tags prime ./cmd
// tinygo flash -target xxx ./cmd

package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/gadget"
	"github.com/merliot/device/runner"
)

var (
	id           = dean.GetEnv("ID", "gadget01")
	name         = dean.GetEnv("NAME", "Gadget")
	deployParams = dean.GetEnv("DEPLOY_PARAMS", "")
	wsScheme     = dean.GetEnv("WS_SCHEME", "ws://")
	port         = dean.GetEnv("PORT", "8000")
	portPrime    = dean.GetEnv("PORT_PRIME", "8001")
	user         = dean.GetEnv("USER", "")
	passwd       = dean.GetEnv("PASSWD", "")
	dialURLs     = dean.GetEnv("DIAL_URLS", "")
	ssids        = dean.GetEnv("WIFI_SSIDS", "")
	passphrases  = dean.GetEnv("WIFI_PASSPHRASES", "")
)

func main() {
	gadget := gadget.New(id, "gadget", name).(*gadget.Gadget)
	gadget.SetDeployParams(deployParams)
	gadget.SetWifiAuth(ssids, passphrases)
	gadget.SetDialURLs(dialURLs)
	gadget.SetWsScheme(wsScheme)
	runner.Run(gadget, port, portPrime, user, passwd, dialURLs, wsScheme)
}
