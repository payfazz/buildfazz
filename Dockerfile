FROM golang:1.11 AS builder

ADD . .

WORKDIR /go/github.com/payfazz/buildfazz

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 $GOPATH/bin/dep
RUN chmod +x $GOPATH/bin/dep
RUN dep init
RUN dep ensure --vendor-only


COPY . ./
RUN go build -o /app cmd/buildfazz/*.go
RUN rm -rf $GOPATH/bin/dep

FROM debian
COPY --from=builder /app .
ENTRYPOINT ["./app"]
