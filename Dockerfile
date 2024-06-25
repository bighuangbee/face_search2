FROM face_jq_dev:latest as builder

WORKDIR /root/face_search

ADD . /root/face_search

#编译环境
ENV PATH=$PATH:/usr/local/go/bin
RUN go env -w GOBIN=/usr/local/go/bin
RUN go env -w GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV GOPATH=/root/go
ENV GOROOT=/usr/local/go

ENV GODEBUG=cgocheck=0
ENV CGO_ENABLED=1

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go env -w GOPRIVATE=git.hiscene.net && \
    git config --global url."git@git.hiscene.net:".insteadOf "https://git.hiscene.net/"

RUN mkdir -p /root/.ssh && chmod 0700 /root/.ssh
ADD id_rsa /root/.ssh/id_rsa
ADD id_rsa.pub /root/.ssh/id_rsa.pub
RUN #chmod 600 /root/.ssh/id_rsa && \
#    chmod 600 /root/.ssh/id_rsa.pub && \
#    ssh-keyscan -t rsa git.hiscene.net >> /root/.ssh/known_hosts
RUN chmod 600 /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa.pub


#git仓库打tag
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /go/cache
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /root/face_search

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/sdk/lib/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/thirdparty/onnxruntime/lib/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/thirdparty/opencv4-ffmpeg/lib/

RUN go mod tidy
RUN cd app/cmd/server && GOOS=linux GOARCH=amd64 go build -o srv-bin .
RUN cd app/cmd/registe && GOOS=linux GOARCH=amd64 go build -o registe-bin .


FROM nvidia/cuda:11.4.3-cudnn8-devel-ubuntu20.04

# Prevent stop building ubuntu at time zone selection.
ENV DEBIAN_FRONTEND=noninteractive


RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys A4B469963BF863CC && \
    sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    apt update


# Build and install ceres solver
RUN apt-get -y install \
    libgoogle-glog-dev \
    libgflags-dev \
    libatlas-base-dev \
    libeigen3-dev \
    libsuitesparse-dev \
    gcc

ENV CUDA_HOME=$CUDA_HOME:/usr/local/cuda
ENV PATH=$PATH:/usr/lib:/usr/local/cuda/bin

WORKDIR /app

COPY --from=builder /root/face_search/libs /app/libs
COPY --from=builder /root/face_search/app/cmd/server/srv-bin /app/srv-bin
COPY --from=builder /root/face_search/app/cmd/registe/registe-bin /app/registe-bin
COPY --from=builder /root/face_search/app/config/config.yaml /app/conf/config.yaml

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/libs/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/libs/sdk/lib/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/libs/thirdparty/onnxruntime/lib/
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/libs/thirdparty/opencv4-ffmpeg/lib/

ENV GODEBUG=cgocheck=0
ENV CGO_ENABLED=1

ENV face_models_path=/app/libs/models

# 暴露22端口
EXPOSE 22
EXPOSE 6002

RUN mkdir logs

CMD ["./srv-bin", "-conf", "/app/conf/config.yaml"]

# docker build -t face_jq_dev .
# docker run -d --runtime=nvidia --gpus all -p 22:22 -p 6002:6002 --privileged=true --name face_jq_dev face_jq_dev
