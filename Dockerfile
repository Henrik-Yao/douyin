FROM golang:latest
ENV GOPROXY https://goproxy.cn
WORKDIR /root/

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/main.go

EXPOSE 8000
ENTRYPOINT ["./main"]