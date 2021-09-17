FROM golang:1.17.1 AS builder
WORKDIR /go/src
RUN go get -d -v firebase.google.com/go
RUN go get -d -v github.com/gin-gonic/gin
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app .
CMD ["./app"]
