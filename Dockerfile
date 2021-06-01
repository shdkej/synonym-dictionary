FROM golang:1.14-alpine
RUN apk add --no-cache git
COPY . .
RUN go build main.go
CMD ./main
