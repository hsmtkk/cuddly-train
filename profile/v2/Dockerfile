FROM golang:1.19 AS builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY env/ /go/src/app/env/
COPY profile/ /go/src/app/profile/
WORKDIR /go/src/app/profile/v2
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /go/src/app/profile/v2/app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]
