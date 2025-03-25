# 构建阶段
FROM golang:1.23.3-alpine AS builder

ENV GO111MODULE=on \
    GOSUMDB=off \
    CGO_ENABLED=0 \
    TZ=Asia/Shanghai

WORKDIR /
COPY . .
RUN go build -ldflags '-w -s' -o /cmd/main ./cmd/main.go

# 运行阶段
FROM alpine:3.19

# 1. 设置运行时环境变量
ENV TZ=Asia/Shanghai \
    APP_ENV=production

# 2. 安装时区数据（如需）
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    apk del tzdata

WORKDIR /
COPY --from=builder /cmd/main /lumino

COPY config/ /config/
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/lumino"]