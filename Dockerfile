FROM golang:1.11-alpine as build

ENV GO111MODULE on
RUN apk add --no-cache git
RUN mkdir -p /go/src/github.com/modoki-paas/modoki

WORKDIR /go/src/github.com/modoki-paas/modoki

COPY . /go/src/github.com/modoki-paas/modoki
RUN CGO_ENABLED=0 go build -o /bin/modoki

FROM scratch
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build /bin/modoki /bin/modoki
COPY --from=build /go/src/github.com/modoki-paas/modoki/swagger /swagger
WORKDIR /
ENTRYPOINT ["/bin/modoki"]
CMD ["--help"]