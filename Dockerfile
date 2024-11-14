FROM ubuntu:20.04
# FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04

COPY . .

RUN apt-get update && apt-get install -y wget
RUN ARCH=$(dpkg --print-architecture) && \
    if [ "$ARCH" = "amd64" ]; then \
        wget -P /tmp http://ftp.cn.debian.org/debian/pool/main/g/golang-1.23/golang-1.23-go_1.23.2-1_amd64.deb && \
        wget -P /tmp http://ftp.cn.debian.org/debian/pool/main/g/golang-1.23/golang-1.23-src_1.23.2-1_all.deb; \
    elif [ "$ARCH" = "arm64" ]; then \
        wget -P /tmp http://ftp.cn.debian.org/debian/pool/main/g/golang-1.23/golang-1.23-go_1.23.2-1_arm64.deb && \
        wget -P /tmp http://ftp.cn.debian.org/debian/pool/main/g/golang-1.23/golang-1.23-src_1.23.2-1_all.deb; \
    else \
        echo "Unsupported architecture: $ARCH" && exit 1; \
    fi
RUN dpkg -i /tmp/golang-1.23-go_*.deb /tmp/golang-1.23-src_*.deb

ENV PATH="/usr/lib/go-1.23/bin:${PATH}"

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct

RUN go build -o sdcs main.go
