FROM ubuntu:20.04
# FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04

ARG GO_VERSION=1.23.2

COPY . .

RUN apt-get update && apt-get install -y wget

RUN ARCH=$(dpkg --print-architecture) && \
    wget -O /tmp/go.tar.gz "https://mirrors.ustc.edu.cn/golang/go${GO_VERSION}.linux-${ARCH}.tar.gz" && \
    tar -xf /tmp/go.tar.gz && \
    rm -f /tmp/go.tar.gz

RUN ./go/bin/go env -w GO111MODULE=on && \
    ./go/bin/go env -w GOPROXY=https://goproxy.cn,direct

RUN ./go/bin/go build -o sdcs
