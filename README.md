
### 环境变量
```
wget https://dl.google.com/go/go1.20.1.linux-amd64.tar.gz
tar -xf go1.20.1.linux-amd64.tar.gz -C /usr/local/
rm -f go1.20.1.linux-amd64.tar.gz


export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/sdk/lib/
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/thirdparty/onnxruntime/lib/
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/root/face_search/libs/thirdparty/opencv4-ffmpeg/lib/

export GODEBUG=cgocheck=0
export CGO_ENABLED=1

export PATH=$PATH:/usr/local/go/bin
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOBIN=/usr/local/go/bin
export GO111MODULE=on
export GOPATH=/root/go
export GOROOT=/usr/local/go

```

### 构建
```
#构建dev
docker build -t face_jq_dev -f Dockerfile.dev .
docker run -d --runtime=nvidia --gpus all -p 22:22 -p 6001:6002 --privileged=true -v C:\hiar_face\registe_path:/hiar_face/registe_path --name face_jq_dev face_jq_dev

构建app
docker build -t face_jq_app .
docker run -d --runtime=nvidia --gpus all -p 23:22 -p 6002:6002 --privileged=true -v C:\hiar_face\registe_path:/hiar_face/registe_path -v C:\hiar_face\logs:/hiar_face/registe_logs  -v C:\hiar_face\search_record:/hiar_face/search_record --env face_models_path=/app/libs/models --name face_jq_app face_jq_app

```

### docker engine
`
{
"builder": {
"gc": {
"defaultKeepStorage": "20GB",
"enabled": true
}
},
"experimental": false,
"features": {
"buildkit": true
},
"insecure-registries": [
"192.168.3.57:82"
],
"registry-mirrors": [
"https://pee6w651.mirror.aliyuncs.com",
"https://registry.docker-cn.com",
"https://docker.mirrors.ustc.edu.cn",
"https://dockerproxy.com",
"https://mirror.ccs.tencentyun.com"
],
"runtimes": {
"nvidia": {
"args": [],
"path": "nvidia-container-runtime"
}
}
}
`
