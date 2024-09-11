

1编译二进制文件   **必须在cmd生成**

```powershell
SET CGO_ENABLED=0  	// 禁用CGO
SET GOOS=linux  	// 目标平台是linux
SET GOARCH=amd64  	// 目标处理器架构是amd64
go build -tags=k8s -o webook-tmh01 .    // 编译成二进制文件
go build -tags=k8s -o webook-live .    // 编译成二进制文件

```

编写dockerfile文件

```dockerfile
#docker 的基础镜像
FROM ubuntu:20.04
# COPY webook-live /app/webook 
COPY webook-tmh01 /app/webook 
WORKDIR /app
CMD ["/app/webook"]
```



在docker中创建server

```powershell
docker rmi -f flycash//webook-tmh01:v0.0.1  // 删除镜像
docker build -t flycash/webook-tmh01:v0.0.1 .  // 构建镜像

docker rmi -f flycash//webook-live:v0.0.1  // 删除镜像
docker build -t flycash/webook-live:v0.0.1 .  // 构建镜像
docker build -t flycash/webook-live:v0.0.2 .
```




### 使用kubectl 来通过镜像生成多个容器

下载kubectl.exe, 然后放到path第一个, 然后执行命令 :

kubectl apply -f k8s-webook-deployment.yaml  	//执行yaml文件

kubectl apply -f k8s-webook-service.yaml



kubectl get deployments   //看当前的deployments

kubectl get pods   //看所有的pod

kubectl get services



kubectl delete deployment webook-live 	//删除webook-live dep...

kubectl delete service webook-live



出现这个错误的处理方式:

Unable to connect to the server: dial tcp: lookup kubernetes.docker.internal: no such host



goland启动用这个验证 [localhost:8080/hello](http://localhost:8080/hello)

k8s启动 验证 get方法  http://localhost/hello    (80端口)

[localhost:82/hello](http://localhost:82/hello)

