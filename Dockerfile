FROM golang:1.9.1

RUN go get -u github.com/kardianos/govendor

WORKDIR /go/src/github.com/impactasaurus/server
COPY vendor/vendor.json vendor/
RUN govendor sync
COPY . .

WORKDIR /go/src/github.com/impactasaurus/server/cmd
RUN go-wrapper install
CMD /go/bin/cmd
