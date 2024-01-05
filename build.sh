# 停止容器
docker stop authsvr
# 删除容器
docker rm authsvr
# 删除镜像
docker rmi authsvr:latest

# 构建新的镜像
go mod tidy
docker build -t authsvr .
docker run -d -p 8849:8849 -p 18849:18849 --name authsvr -v /data/config:/data/config authsvr:latest