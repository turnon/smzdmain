FROM golang:alpine as builder

WORKDIR /go/src
COPY *.go ./

RUN apk add --no-cache git mercurial \
    && git clone https://github.com/golang/net $GOPATH/src/golang.org/x/net \
    && go get github.com/PuerkitoBio/goquery \
    && go get github.com/fatih/color \
    && go build -o smzdm \
    && apk del git mercurial

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /go/src/smzdm .
RUN chmod +x /root/smzdm

ENTRYPOINT ["/root/smzdm"]