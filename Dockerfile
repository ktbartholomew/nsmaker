FROM golang:1.12

WORKDIR /go/src/github.com/ktbartholomew/nsmaker

ADD . .

RUN go build -a -o nsmaker .

CMD ["/go/src/github.com/ktbartholomew/nsmaker/nsmaker"]
