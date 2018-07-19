FROM golang:1.10-alpine as build

RUN apk add --no-cache git
RUN go get -v github.com/cs3238-tsuzu/modoki

WORKDIR /go/src/github.com/cs3238-tsuzu/modoki

COPY . /go/src/github.com/cs3238-tsuzu/modoki
RUN go get -v .
RUN CGO_ENABLED=0 go build -o /bin/modoki

FROM scratch
COPY --from=build /bin/modoki /bin/modoki
COPY --from=build /go/src/github.com/cs3238-tsuzu/modoki/swagger /
WORKDIR /
ENTRYPOINT ["/bin/modoki"]
CMD ["--help"]