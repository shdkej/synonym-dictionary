FROM golang:1.14-alpine
RUN apk add --no-cache git
WORKDIR /go/src/github.com/shdkej/project
COPY . .
CMD go run main.go
