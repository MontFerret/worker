# Worker

**Worker** is a simple HTTP server that accepts FQL queries, executes them and returns their results.
The Docker image contains headless Google Chrome, so feel free to run queries using `cdp` driver.

## Quick start

```sh
docker run -d -p 8080 montferret/worker
```

![worker](https://raw.githubusercontent.com/MontFerret/worker/master/assets/postman.png)
