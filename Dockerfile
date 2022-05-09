ARG IMAGE=golang:1.17-alpine

FROM golang:1.17-alpine AS builder
RUN apk add --no-cache git make
COPY . /go/src/github.com/zikani03/postnat/
RUN cd /go/src/github.com/zikani03/postnat/ && go build -o /dist/postnat cmd/postnat.go

FROM $IMAGE
COPY --from=builder /dist/postnat /bin/postnat
VOLUME /data/
CMD ["/bin/postnat", "--config", "/data/postnat.toml", "run"]
