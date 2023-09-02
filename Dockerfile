# 基础镜像
FROM ubuntu:20.04

WORKDIR /app
# 打包后的go放到镜像的/app/webook目录
COPY webook /app/webook
ENTRYPOINT ["/app/webook"]