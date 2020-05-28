FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
# Make is requiered for build.
RUN apk update && apk add --no-cache git make ca-certificates

WORKDIR /go/src/github.com/MontFerret/worker

COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux make compile

# Build the final container. And install
FROM microbox/chromium-headless:83.0.4103.0 as runner

RUN apt-get update && apt-get install -y dumb-init

WORKDIR /root

# Add in certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.c

# Add worker binary
COPY --from=builder /go/src/github.com/MontFerret/worker/bin/worker .

EXPOSE 8080

ENTRYPOINT ["dumb-init", "--"]
CMD ["/bin/sh", "-c", "chromium --no-sandbox --disable-setuid-sandbox --disable-gpu --headless --remote-debugging-port=9222 & ./worker"]