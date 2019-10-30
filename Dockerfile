FROM golang:1.12 as build

WORKDIR $GOPATH/src/github.com/BardoBravo/real-estate-finder
COPY scraper scraper
COPY main.go .

RUN go get -d -v ./...
RUN go install

FROM gcr.io/distroless/base

COPY --from=build /go/bin/real-estate-finder /
CMD ["/real-estate-finder"]
