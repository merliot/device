#### Install Instructions

1. Click the download link to obtain the installer file for your device.  The installer is built to run on Raspberry Pi OS and will install the device as a service with automatic restart on boot.

2. (Optional) Verify the MD5 checksum matches the installer file MD5 checksum.

    ```
    $ md5sum installer
    ```

3. Make the installer executable:

    ```
    $ chmod +x installer
    ```

4. Run the installer:

    ```
    $ sudo ./installer
    ```

#### Uninstall

To uninstall, use the -u option:

```
    $ sudo ./installer -u
```
