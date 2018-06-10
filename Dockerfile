FROM golang:1.10-alpine as build

RUN apk add --no-cache git

WORKDIR /go/src/github.com/cs3238-tsuzu/modoki

COPY . /go/src/github.com/cs3238-tsuzu/modoki
RUN go get -v .
RUN CGO_ENABLED=0 go build -o /bin/modoki

FROM scratch
COPY --from=build /bin/modoki /bin/modoki
ENTRYPOINT ["/bin/modoki"]
CMD ["--help"]