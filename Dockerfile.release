# Build the final container. And install
FROM microbox/chromium-headless:75.0.3765.1 as runner

RUN apt-get update && apt-get install -y dumb-init

WORKDIR /root

# Add worker binary
COPY worker /bin/worker
EXPOSE 8080

ENTRYPOINT ["dumb-init", "--"]
CMD ["/bin/sh", "-c", "chromium --no-sandbox --disable-setuid-sandbox --disable-gpu --headless --remote-debugging-port=9222 & /bin/worker"]