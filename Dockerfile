FROM golang:1.13-alpine3.10 as builder

ENV CGO_ENABLED 0
ENV GO111MODULE on
RUN apk --no-cache add git
WORKDIR /go/src/kubethanos
COPY . .
RUN go run main.go --namespaces=!kube-system
ENV GOARCH amd64
RUN go build -o /bin/kubethanos -v \
  -ldflags "-X main.version=$(git describe --tags --always --dirty) -w -s"

FROM alpine:3.10
MAINTAINER Berkay Dinçer <dincerbberkay@gmail.com>

RUN apk --no-cache add ca-certificates dumb-init tzdata
COPY --from=builder /bin/kubethanos /bin/kubethanos

USER 65534
ENTRYPOINT ["/bin/kubethanos"]
