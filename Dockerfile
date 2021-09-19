FROM golang:1.16.6-buster as builder
WORKDIR /app
ARG VERSION=dev
ARG GOFLAGS
COPY . /app
RUN make build

FROM debian:10.9-slim as final
RUN set -ex &&\
 apt-get update &&\
 DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates tzdata &&\
 rm -rf /var/lib/apt/lists/*
WORKDIR /data/Bukukas-Inmemory-Search
EXPOSE 8080
ENTRYPOINT ["/data/Bukukas-Inmemory-Search/Bukukas-Inmemory-Search"]
COPY . /app
COPY --from=builder /app/build/* /data/Bukukas-Inmemory-Search/
