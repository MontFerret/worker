# Worker

**Worker** is the simple HTTP server that accepts FQL queries, execute it and return the result.

## Docker

Image contains headless Google Chrome so feel free to run queries using `cdp` driver.

1. Build image
```sh
docker build -t ferret-worker .
```

2. Run container
```
docker run -p 8080:8080 -it ferret-worker
```