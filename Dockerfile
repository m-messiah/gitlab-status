FROM alpine:latest
RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*
EXPOSE 8080
ARG server_url
ARG gitlab_url
ARG app_key
ARG app_secret
ENV REDIRECT_URL=https://${server_url}/oauth-authorized \
	GITLAB_URL=${gitlab_url} \
	GITLAB_APP_KEY=${app_key} \
	GITLAB_APP_SECRET=${app_secret}
ADD status /status
RUN chmod +x /status
ENTRYPOINT /status
