FROM golang:alpine AS builder
ENV GOPROXY=https://goproxy.cn
ENV TZ Asia/Shanghai

#RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
#    && echo ${TZ} > /etc/timezone \
#    && apk del tzdata

WORKDIR /app

ADD go.mod .
COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o timer service/cmd/foundation/timer/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/timer /app/timer

CMD ["./timer"]