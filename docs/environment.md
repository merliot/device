## Device Environment Variables

These environment variables configure the device:

**ID**

Device ID.

**NAME**

Device name.

**WS_SCHEME**

WebSocket scheme to use when connecting to device.  Defaults to ws://.  Set to wss:// for device running under https.

**PORT**

Port the device listens on, default is 8000.

**PORT_PRIME**

Port the prime device listens on, default is 8001.

**USER, PASSWD**

Set user and password for HTTP Basic Authentication on the device.  The user will be prompted for user/password when browsing to the device.

**DIAL_URLS**

Additional URLs for the device to dial.  The device will always dial to the host it was deployed on.  Set DIAL_URLS to additional hosts for the device to dial into.

**WIFI_SSIDS, WIFI_PASSPHRASES**

Set Wifi SSID(s) and passphrase(s) for Wifi-enabled devices built with TinyGo.  These are matched comma-delimited lists.  For each SSID, there should be a matching passphrase.  For example:

- WIFI_SSIDS="test,backup"
- PASSPHRASES="testtest,ihavenoplan"

So testtest goes with SSID test, and ihavenoplan goes with SSID backup.
