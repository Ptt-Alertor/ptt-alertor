# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM alpine:latest

RUN set -eux; \
	apk add --no-cache --virtual ca-certificates

# Copy the local package files to the container's workspace.
ADD ptt-alertor /
COPY config/ config/
COPY public/ public/

# Run the outyet command by default when the container starts.
ENTRYPOINT /ptt-alertor

# Document that the service listens on port 9090.
EXPOSE 9090 6060