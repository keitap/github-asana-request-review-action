FROM golang:1.16-alpine

RUN apk add --no-cache git

WORKDIR /go/src/app
COPY . .
RUN go build -o action ./cmd
ENTRYPOINT ["/go/src/app/action"]
