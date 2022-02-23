FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go build -o /status

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /status /status
EXPOSE 8080
ENTRYPOINT /status
