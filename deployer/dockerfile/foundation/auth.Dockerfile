FROM golang:alpine AS builder
ENV GOPROXY="https://goproxy.cn,direct"
ENV TZ Asia/Shanghai

#RUN apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
#    && echo ${TZ} > /etc/timezone \
#    && apk del tzdata

WORKDIR /app

ADD go.mod .
COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth service/cmd/foundation/auth/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/auth /app/auth
COPY --from=builder /app/service/cmd/foundation/auth/rbac_model.conf /app/rbac_model.conf

CMD ["./auth"]