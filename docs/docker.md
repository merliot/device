# Run Device in a Docker Container

A device can be run in a docker container.  The device has a Dockerfile used to build a docker image.

### 1. Git clone device (using Merliot Hub in this example):

```
git clone https://github.com/merliot/hub.git
cd hub
// see Dockerfile
```

### 2. Build docker image for device:

The device can be built to listen over http:// or https://

#### Build for https://

```
sudo docker build -t hub --build-arg .
```

#### Build for http://

```
sudo docker build -t hub --build-arg SCHEME=ws .
```

### 3. Run docker container for device:

```
sudo docker run -p 8000:8000 hub
```

Browse to http://localhost:8000

