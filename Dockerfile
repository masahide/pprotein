# --------------------------------------------------

FROM golang:1.23-alpine AS pprotein

RUN apk add npm make

WORKDIR $GOPATH/src/app
COPY . .

RUN make build

# --------------------------------------------------

FROM golang:1.23-alpine AS alp

RUN go install github.com/tkuchiki/alp/cmd/alp@latest

# --------------------------------------------------

FROM golang:1.23-alpine AS slp

RUN apk add gcc musl-dev
RUN go install github.com/tkuchiki/slp/cmd/slp@latest

# --------------------------------------------------

FROM alpine as percona

RUN apk add --no-cache tar curl
RUN curl https://downloads.percona.com/downloads/percona-toolkit/3.6.0/binary/tarball/percona-toolkit-3.6.0_x86_64.tar.gz -o /tmp/percona-toolkit.tar.gz
RUN tar -xzf /tmp/percona-toolkit.tar.gz -C /tmp

# --------------------------------------------------

FROM alpine

RUN apk add --no-cache graphviz \
  perl # for pt-query-digest

COPY --from=pprotein /go/src/app/pprotein /usr/local/bin/
COPY --from=pprotein /go/src/app/pprotein-agent /usr/local/bin/
COPY --from=alp /go/bin/alp /usr/local/bin/
COPY --from=slp /go/bin/slp /usr/local/bin/
COPY --from=percona /tmp/percona-toolkit-3.6.0/bin/ /usr/local/bin/

RUN mkdir -p /opt/pprotein
WORKDIR /opt/pprotein

ENTRYPOINT ["pprotein"]
