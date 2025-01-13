FROM golang:1.23.4-alpine

RUN apk add --no-cache git
ADD . /go/src

WORKDIR /go/src

RUN go mod tidy
RUN go install

ENV PATH="/go/bin:$PATH"
