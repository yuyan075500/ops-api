FROM golang:1.21.8

# 指定工作目录
WORKDIR /app

# 复制源到工作目录
COPY . /app

# Go环境设置
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

# 处理项目依赖
RUN go mod tidy

# 项目编译
RUN go build

# 暴露端口
EXPOSE 8000

# 执行
ENTRYPOINT ["./ops-api"]
