#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/playersvc -v ./playersvc/cmd/playersvc/main.go

#final stage
FROM alpine:latest
COPY --from=builder /go/bin/playersvc /playersvc
ENTRYPOINT /playersvc
LABEL Name=playersvc Version=0.0.1
EXPOSE 8085
