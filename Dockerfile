### Builder for Golang
FROM golang:1.19 as go-builder
MAINTAINER @jedipunkz

WORKDIR /go/src/
ADD . /go/src/

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/esp

## Deployable image
FROM alpine:latest as esp-base
MAINTAINER @jedipunkz

COPY --from=go-builder /go/bin/esp /go/bin/esp

ENTRYPOINT ["/bin/sh",  "-c" , "/go/bin/esp"]
