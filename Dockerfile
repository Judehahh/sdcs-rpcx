FROM ubuntu:20.04
# FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04

COPY . .

RUN apt-get update && apt-get install -y wget
RUN wget -P /tmp http://ftp.cn.debian.org/debian/pool/main/g/golang-1.23/golang-1.23-go_1.23.2-1_amd64.deb \
    && wget -P /tmp http://ftp.cn.debian.org/debian/pool/main/g/golang-1.23/golang-1.23-src_1.23.2-1_all.deb
RUN dpkg -i /tmp/golang-1.23-go_1.23.2-1_amd64.deb /tmp/golang-1.23-src_1.23.2-1_all.deb

ENV PATH="/usr/lib/go-1.23/bin:${PATH}"

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct

RUN go build -o sdcs main.go
