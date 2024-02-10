// go run ./cmd
// go run -tags prime ./cmd
// tinygo flash -target xxx ./cmd

package main

import (
	"github.com/merliot/dean"
	"github.com/merliot/device/runner"
	"github.com/merliot/skeleton"
)

var (
	id           = dean.GetEnv("ID", "skeleton01")
	name         = dean.GetEnv("NAME", "Skeleton")
	deployParams = dean.GetEnv("DEPLOY_PARAMS", "target=demo")
	port         = dean.GetEnv("PORT", "8000")
	portPrime    = dean.GetEnv("PORT_PRIME", "8001")
	user         = dean.GetEnv("USER", "")
	passwd       = dean.GetEnv("PASSWD", "")
	dialURLs     = dean.GetEnv("DIAL_URLS", "")
	ssids        = dean.GetEnv("WIFI_SSIDS", "")
	passphrases  = dean.GetEnv("WIFI_PASSPHRASES", "")
)

func main() {
	s := skeleton.New(id, "skeleton", name).(*skeleton.Skeleton)
	s.SetDeployParams(deployParams)
	s.SetWifiAuth(ssids, passphrases)
	s.SetDialURLs(dialURLs)
	runner.Run(s, port, portPrime, user, passwd, dialURLs)
}
