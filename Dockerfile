FROM golang:1.14.0 as build

WORKDIR /go/src/github.com/dtamura/hello-mongo
COPY . .
RUN make

FROM alpine:3.11.5
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /go/src/github.com/dtamura/hello-mongo/bin/hello-mongo app

CMD [ "./app", "start" ]