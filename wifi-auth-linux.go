//go:build !tinygo

package device

import (
	"os"
	"strconv"
)

func (d *Device) ParseWifiAuth() {
	if ssid, ok := os.LookupEnv("WIFI_SSID"); ok {
		if passphrase, ok := os.LookupEnv("WIFI_PASSPHRASE"); ok {
			d.WifiAuth[ssid] = passphrase
		}
	}
	for i := 0; i < 10; i++ {
		a := strconv.Itoa(i)
		if ssid, ok := os.LookupEnv("WIFI_SSID_" + a); ok {
			if passphrase, ok := os.LookupEnv("WIFI_PASSPHRASE_" + a); ok {
				d.WifiAuth[ssid] = passphrase
			}
		}
	}
}
