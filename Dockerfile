FROM golang:1.14.0

WORKDIR /go/src/github.com/dtamura/hello-mongo
COPY . .
RUN go get -d -v .
RUN go build

ENTRYPOINT [ "./hello-mongo" ]