# syntax=docker/dockerfile:1

FROM golang

WORKDIR /app

ADD . /app/

RUN go version

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN go build -o backend .

EXPOSE 7013

CMD /app/backend

