FROM golang:1.21-bullseye
EXPOSE 5678
WORKDIR /workspace
RUN sed -i s/deb.debian.org/mirrors.aliyun.com/g /etc/apt/sources.list \
        && sed -i s/security.debian.org/mirrors.aliyun.com/g /etc/apt/sources.list \
        && apt-get update -y&& apt-get upgrade -y  && apt-get install -y make


# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
RUN go mod tidy
# src code
COPY . .
RUN make ugrsicalsrv-linux-amd64
RUN mv build/ugrsicalsrv-linux-amd64 ./ugrsicalsrv
RUN chmod +x ./ugrsicalsrv

ENV TZ=Asia/Shanghai
RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
        && echo ${TZ} > /etc/timezone
ENTRYPOINT ["./ugrsicalsrv"]