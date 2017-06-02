FROM golang:1.8

WORKDIR /go/src/github.com/impactasaurus/server/cmd
COPY . /go/src/github.com/impactasaurus/server

RUN go-wrapper download
RUN go-wrapper install

CMD /go/bin/cmd
