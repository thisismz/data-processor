ARG GO_VERSION=1.22.2
ARG BASH_VERSION=5.0


FROM golang:$GO_VERSION AS goBuilder
COPY ./ /data-processor

RUN cd /data-processor \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=Stable -X main.Build=$(date -u '+%Y-%m-%d_%H:%M:%S')"  -buildvcs=false -o data-processor  ./cmd/main.go 



FROM amd64/bash:$BASH_VERSION
RUN sed -i -e 's/http:/https:/' /etc/apk/repositories \
    && apk --no-cache add busybox-suid curl rsync tzdata tcpdump ca-certificates \
    && mkdir -p /app/log

WORKDIR /go/bin/

COPY --from=goBuilder /data-processor/data-processor .


# Set the binary.
ENTRYPOINT ["/go/bin/data-processor"]