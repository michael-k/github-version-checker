FROM golang:1.10.3-alpine3.8 as builder

RUN apk add --no-cache git

COPY . /go/src/github.com/michael-k/github-version-checker

RUN go get -v ./...

RUN cd /go/src/github.com/michael-k/github-version-checker/ && go build

# go get github.com/mcuadros/go-version
# go get golang.org/x/oauth2
# go get github.com/shurcooL/githubv4

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/github-version-checker /github-version-checker
ENTRYPOINT ["/github-version-checker"]
