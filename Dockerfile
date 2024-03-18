# syntax=docker/dockerfile:1

# Debian GNU/Linux 12 (bookworm)
FROM golang:1.22.0

RUN wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.1/tinygo_0.31.1_amd64.deb
RUN dpkg -i tinygo_0.31.1_amd64.deb

RUN apt-get update
RUN apt-get install vim tree bc -y
