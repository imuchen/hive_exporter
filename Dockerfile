FROM golang:1.10-alpine as builder
RUN mkdir -p $GOPATH/src/app
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/app/app .
CMD ["./app"]
