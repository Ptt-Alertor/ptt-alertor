FROM buildpack-deps:jessie-scm
# gcc for cgo

RUN apt-get update && apt-get install -y --no-install-recommends \
   g++ \
   gcc \
   libc6-dev \
   make \
   pkg-config \
   && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.8
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 53ab94104ee3923e228a2cb2116e5e462ad3ebaeea06ff04463479d7f12d27ca

RUN curl -kfsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go/
ENV GO_WORKDIR github.com/liam-lai/ptt-alertor/
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY docker_golang_1.8/go-wrapper /usr/local/bin/

ADD . "$GOPATH"/src/"$GO_WORKDIR"

RUN go get "$GO_WORKDIR"
RUN go install "$GO_WORKDIR"

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/ptt-alertor

# Document that the service listens on port 9090.
EXPOSE 9090
