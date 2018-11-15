FROM golang:1.10.3-alpine

WORKDIR /go/src/app
COPY . .

RUN apk update && \
    apk upgrade && \
    apk add git

RUN go get -v ./...

ENTRYPOINT ["app"]