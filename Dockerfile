# GO Builder
###############
FROM golang:alpine AS goBuilder

RUN apk update && apk add --no-cache git make ca-certificates
WORKDIR /go/src/github.com/MontFerret/worker
COPY . .

RUN CGO_ENABLED=0 GOOS=linux make compile

# MITM Builder
###############
FROM pierrebrisorgueil/mitm:latest AS mitmBuilder

# Runner
###############
FROM montferret/chromium:91.0.4469.0 as runner
RUN apt-get update && apt-get install -y dumb-init

# mitm
RUN apt-get update && apt-get install --no-install-recommends -y python3.8 python3-pip python3.8-dev
RUN pip install mitmproxy bs4 lxml
COPY --from=mitmBuilder bundle.js /
COPY --from=mitmBuilder inject.py /

# worker
COPY --from=goBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.c
COPY --from=goBuilder /go/src/github.com/MontFerret/worker/bin/worker .

EXPOSE 8080

ENTRYPOINT ["dumb-init", "--"]
CMD ["/bin/sh", "-c", "mitmdump -p 8081 -s inject.py & CHROME_OPTS='--proxy-server=127.0.0.1:8081' ./entrypoint.sh & ./worker"]