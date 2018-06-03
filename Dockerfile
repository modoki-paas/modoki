FROM golang:1.10-alpine as build

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/cs3238-tsuzu/modoki/
WORKDIR /go/src/github.com/cs3238-tsuzu/modoki

COPY . /go/src/github.com/cs3238-tsuzu/modoki
RUN go get -v .
RUN go build -o /bin/modoki

FROM scratch
COPY --from=build /bin/modoki /bin/modoki
ENTRYPOINT ["/bin/modoki"]
CMD ["--help"]