# GitLab Build Status

Simple dashboard with Gitlab oauth.

+ Shows only available projects, with enabled CI.
+ Shows each job status separately
+ Shows coverage for project, if exists.

![screenshot](screenshot.png)

## Run

The project initially was intended to run as docker service, but could be easily converted to anything else (k8s, [systemd](/status.service.j2), sysv).

The easies way is to:
1. Build the docker image from [Dockerfile](/Dockerfile) provided
2. Run the container with environment variables, where `${server_url}` will be external address of this container:
	```
	REDIRECT_URL=https://${server_url}/oauth-authorized
	GITLAB_URL=${gitlab_url}
	GITLAB_APP_KEY=${app_key}
	GITLAB_APP_SECRET=${app_secret}
	```
3. Expose port 8080 of the port through some ingress or reverse-proxy with the name `${server_url}`
