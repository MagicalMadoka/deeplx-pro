FROM golang:1.21 AS builder
WORKDIR /go/src/github.com/OwO-Network/DeepLX
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY login.go ./
RUN go get -d -v ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o deeplx .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/OwO-Network/DeepLX/deeplx /app/deeplx
CMD ["/app/deeplx"]
