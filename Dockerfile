FROM alpine:latest
RUN apk update && apk add ca-certificates
ADD status /status
RUN chmod +x /status
ENTRYPOINT /status
