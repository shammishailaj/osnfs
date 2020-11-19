FROM golang:alpine

RUN mkdir /go/src/osnfs

WORKDIR /go/src/osnfs

COPY . ./

RUN go env

RUN go build -o ./cmd/osnfs/osnfs ./cmd/osnfs/osnfs.go

CMD ./cmd/osnfs/osnfs