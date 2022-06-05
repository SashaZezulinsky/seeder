FROM golang:alpine as builder
RUN apk add make
ADD . /seeder
RUN cd /seeder && make build

FROM alpine:3.14
COPY --from=builder /seeder/bin/* /usr/local/bin/
