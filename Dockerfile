FROM golang:1.11.2-alpine3.8 as builder

RUN apk add --no-cache git

COPY . /go/src/github.com/michael-k/github-version-checker
WORKDIR /go/src/github.com/michael-k/github-version-checker

RUN go get -v -x ./...
RUN go build


FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/github-version-checker /github-version-checker
ENTRYPOINT ["/github-version-checker"]
