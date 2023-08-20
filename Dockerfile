# build start
FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -ldflags="-s -w" -o /app/main main.go
# build end

# run start
FROM alpine:3.17

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai


COPY --from=builder /app/main /qnhd/main
WORKDIR /qnhd
COPY boring-avatars-service /qnhd/avatar
COPY Docker/conf /qnhd/conf
COPY Docker/dict /qnhd/dict
COPY Docker/run.sh run.sh
RUN mkdir -p /qnhd/runtime

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata && \
    apk add ca-certificates curl nodejs npm

RUN cd /qnhd/avatar && npm install && cd /qnhd

EXPOSE 7013

RUN chmod +x /qnhd/main && chmod +x /qnhd/run.sh
ENV RELEASE=1
CMD [ "sh", "/qnhd/run.sh" ]

# run end
