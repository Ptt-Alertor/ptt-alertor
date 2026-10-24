# building binary
FROM golang:1.27-alpine as builder

ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/ptt-alertor .

# building executable image
FROM alpine:latest

RUN set -eux; \
	apk add --no-cache ca-certificates

COPY public/ public/
COPY --from=builder /out/ptt-alertor .

ENTRYPOINT /ptt-alertor
