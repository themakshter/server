FROM golang:1.8

RUN go get -u github.com/kardianos/govendor

WORKDIR /go/src/github.com/impactasaurus/server
COPY . .
RUN govendor get

WORKDIR /go/src/github.com/impactasaurus/server/cmd
RUN go-wrapper install
CMD /go/bin/cmd
