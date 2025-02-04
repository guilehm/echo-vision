FROM golang:1.23.5-alpine

RUN apk add --no-cache git curl

WORKDIR /go/src
ADD . /go/src

ENV GOWORK=/go/src/go.work
RUN go work sync

WORKDIR /go/src/echo-analyzer

RUN curl -sSfL https://release.ariga.io/atlas/atlas-linux-amd64-latest -o /usr/local/bin/atlas \
    && chmod +x /usr/local/bin/atlas

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

ENV PATH="/go/bin:$PATH"
