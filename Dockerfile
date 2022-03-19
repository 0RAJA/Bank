# 导入官方镜像
FROM golang:1.16-alpine AS builder

#为镜像设置环境变量
ENV Go111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN go build -o cmd/bank/main cmd/bank/bank.go
#RUN apk add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine

WORKDIR /app

#COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/cmd/bank/main .
COPY --from=builder /app/configs/config.yml ./configs/

EXPOSE 8080

CMD ["/app/main"]
