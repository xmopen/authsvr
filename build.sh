# 停止容器
docker stop authsvr
# 删除容器
docker rm authsvr
# 删除镜像
docker rmi authsvr:latest

# 构建新的镜像
docker build -t authsvr .

go mod tidy
go build -o ./svrmain ./*.go
# 限制内存为20M，CPU使用核数为0.1核

docker run -d -p 8849:8849 -p 18849:18849 -e TZ=Asia/Shanghai --memory=20m --cpus=0.1 --oom-kill-disable=true --name authsvr -v /data/config:/data/config authsvr:latest