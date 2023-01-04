FROM golang:latest as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/xihe-finetune
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o xihe-finetune .

# copy binary config and utils
FROM alpine:3.14
RUN apk update && apk add --no-cache \
        git \
        bash \
        libc6-compat
COPY --from=BUILDER /go/src/github.com/opensourceways/xihe-finetune/xihe-finetune /opt/app/xihe-finetune

ENTRYPOINT ["/opt/app/xihe-finetune"]
