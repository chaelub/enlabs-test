FROM golang:latest

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app
RUN go mod tidy
RUN mkdir bin
RUN go build -o bin/enlabs main.go
CMD ["/go/src/app/bin/enlabs"]