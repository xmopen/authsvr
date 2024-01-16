FROM centos
LABEL authors="openxm"

ENV TZ=Asia/Shanghai

# 拷贝配置文件、以及二进制脚本
COPY ./svrmain /data/code/go/authsvr/
COPY ./config/ /data/config/

# HTTP: 8849
# RPC : 18849
EXPOSE 8849 18849
ENTRYPOINT ["/data/code/go/authsvr/svrmain"]
