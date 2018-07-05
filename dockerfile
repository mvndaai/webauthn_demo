FROM golang:1.10.3-alpine

WORKDIR /go/src/app
COPY . .

RUN apk update && \
    apk upgrade && \
    apk add git

RUN go get -d github.com/nanobox-io/golang-scribble
RUN go get -d github.com/labstack/echo

RUN go install -v ./...
CMD ["app"]