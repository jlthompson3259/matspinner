#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/ticketsvc -v ./cmd/ticketsvc/main.go

#final stage
FROM alpine:latest
COPY --from=builder /go/bin/ticketsvc /ticketsvc
ENTRYPOINT /ticketsvc
LABEL Name=ticketsvc Version=0.0.1
EXPOSE 8085
EXPOSE 8086
