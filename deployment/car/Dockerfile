# 启动编译环境
FROM golang:1.16-alpine AS buidler
# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
# 拷贝源代码到镜像中
COPY . /go/src/coolcar/server
# 编译
WORKDIR /go/src/coolcar/server
RUN go install ./car/...

FROM alpine:3.13
COPY --from=buidler /go/bin/car /bin/car
ENV ADDR=8081
ENV WS_ADDR=:8080

#暴露端口  运行时通过-p进行映射
EXPOSE 8080
EXPOSE 8081
# 设置服务入口 ENTRYPOINT写死的  CMD是默认参数 CMD [ "executable" ]可以让用户重写
ENTRYPOINT [ "/bin/car" ]
