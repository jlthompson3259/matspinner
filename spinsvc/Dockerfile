#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/spinsvc -v ./spinsvc/cmd/spinsvc/main.go

#final stage
FROM alpine:latest
COPY --from=builder /go/bin/spinsvc /spinsvc
ENTRYPOINT /spinsvc
LABEL Name=spinsvc Version=0.0.1
EXPOSE 8086
