FROM golang:1.20
LABEL authors="zhenxinma"
# 1、设置工作环境.
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"
# 2、在容器内设置/data/app为当前工作目录.
# WORKDIR 不存在则会创建.
WORKDIR /data/app/bin

# 3、将当前文件复制到工作目录以及配置文件
COPY . .
# 拷贝配置文件
RUN mkdir /data/config

# 4、拷贝svrmain
COPY ./svrmain .

# 5、暴露端口
# http 端口
# rpc 端口
EXPOSE 8849 18849

# 6、容器入口点
ENTRYPOINT ["/data/app/bin/svrmain"]