FROM golang:alpine AS builder
ENV GOPROXY=https://goproxy.cn
ENV TZ Asia/Shanghai

#RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
#    && echo ${TZ} > /etc/timezone \
#    && apk del tzdata

WORKDIR /app

ADD go.mod .
COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o engine service/cmd/workflow-app/engine/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/engine /app/engine

CMD ["./engine"]