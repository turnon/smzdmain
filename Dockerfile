FROM golang:1.13.3-alpine3.10 as builder

WORKDIR $GOPATH/src/github.com/turnon/smzdm
COPY . ./

RUN apk add --no-cache git \
    && export GO111MODULE=on GOPROXY=https://goproxy.io \
    && go get ./... \
    && go build -o /smzdm \
    && apk del git

FROM alpine:3.10

WORKDIR /root/
COPY --from=builder /smzdm .
RUN chmod +x /root/smzdm

ENTRYPOINT ["/root/smzdm"]