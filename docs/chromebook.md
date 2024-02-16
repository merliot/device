Install Go

```
// See latest version at https://go.dev

wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo  rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

// add to .bashrc

export PATH=$PATH:/usr/local/go/bin
```

Create work dir

```
mkdir work
cd work
```

Git clone device (using Merliot Hub in this example):

```
https://github.com/merliot/hub.git
```

Run device locally:

```
go run ./cmd
```

Browse to http://localhost:8000.
