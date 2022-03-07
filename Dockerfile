FROM golang:1.17 as builder

WORKDIR /go/src/app
COPY . .
RUN go build -o action ./cmd

FROM gcr.io/distroless/base-debian11

LABEL org.opencontainers.image.source https://github.com/keitap/github-asana-request-review-action

COPY --from=builder /go/src/app/action /action
ENTRYPOINT ["/action"]
