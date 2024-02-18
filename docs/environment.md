## Device Environment Variables

These variables configure the hub and devices:

**PORT**

Port the hub listens on, default is 8000.

**USER, PASSWD**

Set user and password for HTTP Basic Authentication on the hub.  The user will be prompted for user/password when browsing to the hub.  These values (if set) are automatically passed down to the device when deployed, and the device connects to the hub using these creditials.

**WIFI_SSIDS, WIFI_PASSPHRASES**

Set Wifi SSID(s) and passphrase(s) for Wifi-enabled devices built with TinyGo.  These are matched comma-delimited lists.  For each SSID, there should be a matching passphrase.  For example:

- WIFI_SSIDS="test,backup"
- PASSPHRASES="testtest,ihavenoplan"

So testtest goes with SSID test, and ihavenoplan goes with SSID backup.
