# Run Device on a Chromebook using Docker

Run a device on a Chromebook using Docker.

### 1. Install [Linux](https://support.google.com/chromebook/answer/9145439?hl=en) on Chromebook.

### 2. Install [Docker](https://www.docker.com/)

These [instructions](https://dvillalobos.github.io/2020/How-to-install-and-run-Docker-on-a-Chromebook/) worked good for me.

### 3. Git clone device (using Merliot Hub in this example):

```
git clone https://github.com/merliot/hub.git
```

### 4. Build docker image for device:

```
cd hub
sudo docker build -t hub --build-arg SCHEME=ws .
```

### 5. Run docker container for device:

```
sudo docker run -p 8000:8000 hub
```

Browse to http://localhost:8000
