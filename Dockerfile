FROM golang:1.14.0 as build

WORKDIR /go/src/github.com/dtamura/hello-mongo
COPY . .
RUN go get -d -v  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app

FROM alpine:3.11.5
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /go/src/github.com/dtamura/hello-mongo/app .

CMD [ "./app" ]