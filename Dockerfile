FROM golang:1.15.7 as builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go  ./

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/datadog-github-pr \
    -ldflags '-s -w'

FROM alpine:3.12.3 as runner

COPY --from=builder /go/bin/datadog-github-pr /app/datadog-github-pr

ENTRYPOINT ["/app/datadog-github-pr"]
