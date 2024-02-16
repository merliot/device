# Development on a Chromebook

Use a Chromebook to develop devices.

### Install [Linux](https://support.google.com/chromebook/answer/9145439?hl=en) on Chromebook.

### Install [Go](https://go.dev):

```
// Get latest version at https://go.dev

wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo  rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

// add to .bashrc

export PATH=$PATH:/usr/local/go/bin
```

### Create work dir

```
mkdir work
cd work
```

### Git clone device (using Merliot Hub in this example):

```
git clone https://github.com/merliot/hub.git
cd hub
go work use .
```

### Run device locally:

```
go run ./cmd
```

Browse to http://localhost:8000.
