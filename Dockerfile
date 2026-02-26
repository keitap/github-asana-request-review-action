FROM golang:1.25 AS builder

WORKDIR /go/src/app
COPY . .
RUN go build -o action ./cmd

FROM gcr.io/distroless/base-debian12

LABEL org.opencontainers.image.source https://github.com/keitap/github-asana-request-review-action

COPY --from=builder /go/src/app/action /action
ENTRYPOINT ["/action"]
