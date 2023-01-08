### Builder for Golang
FROM golang:1.19 as go-builder
MAINTAINER @jedipunkz

WORKDIR /go/src/
ADD . /go/src/

RUN go mod download
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/esp
RUN CGO_ENABLED=0 go build -o /go/bin/esp

# ENTRYPOINT ["/go/bin/esp"]
## Deployable image
FROM alpine:latest as esp-base
MAINTAINER @jedipunkz

COPY --from=go-builder /go/bin/esp /go/bin/esp
# COPY --from=go-builder /go/src/entrypoint.sh /go/bin/entrypoint.sh

# RUN /go/bin/esp
# CMD ["/usr/bin/tini", "--", "/go/bin/entrypoint.sh"]
# ENTRYPOINT ["/go/bin/esp"]
# ENTRYPOINT ["/bin/sh", "-c", "while :; do sleep 10; done"]
# ENTRYPOINT ["/bin/sh", "/go/bin/esp"]
ENTRYPOINT ["/bin/sh",  "-c" , "/go/bin/esp"]
