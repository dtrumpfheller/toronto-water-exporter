#########################
# Build
#########################

FROM golang:1.19 as builder

WORKDIR /go/src/github.com/dtrumpfheller/toronto-water-exporter

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

COPY *.go .
COPY helpers/*.go ./helpers/
COPY torontowater/*.go ./torontowater/
COPY influxdb/*.go ./influxdb/

RUN CGO_ENABLED=0 go build -o /go/bin/app .


#########################
# Deploy
#########################

FROM gcr.io/distroless/static

USER nobody:nobody

COPY --from=builder --chown=nobody:nobody /go/bin/app /toronto-water-exporter/

ENTRYPOINT ["/toronto-water-exporter/app"]