FROM golang:1.23.5-alpine

RUN apk add --no-cache git
ADD . /go/src

WORKDIR /go/src

RUN go mod tidy
RUN go install

RUN apk add --no-cache curl \
    && curl -sSfL https://release.ariga.io/atlas/atlas-linux-amd64-latest -o /usr/local/bin/atlas \
    && chmod +x /usr/local/bin/atlas

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

ENV PATH="/go/bin:$PATH"
