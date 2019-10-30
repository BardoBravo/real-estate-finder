FROM golang:1.12 as build

WORKDIR $GOPATH/src/github.com/BardoBravo/finder-imob
COPY scraper scraper
COPY main.go .

RUN go get -d -v ./...
RUN go install

FROM gcr.io/distroless/base

COPY --from=build /go/bin/finder-imob /
CMD ["/finder-imob"]
