FROM golang:1.14-alpine
RUN apk add --no-cache git
WORKDIR /go/src/github.com/shdkej/project
COPY . .
ENV ELASTICSEARCH_HOST=synonym-es
CMD go run main.go
