# Worker

<p align="center">
	<a href="https://goreportcard.com/report/github.com/MontFerret/worker">
		<img alt="Go Report Status" src="https://goreportcard.com/badge/github.com/MontFerret/worker">
	</a>
<!-- 	<a href="https://codecov.io/gh/MontFerret/worker">
		<img alt="Code coverage" src="https://codecov.io/gh/MontFerret/worker/branch/master/graph/badge.svg" />
	</a> -->
	<a href="https://discord.gg/kzet32U">
		<img alt="Discord Chat" src="https://img.shields.io/discord/501533080880676864.svg">
	</a>
	<a href="https://github.com/MontFerret/worker/releases">
		<img alt="Lab release" src="https://img.shields.io/github/release/MontFerret/worker.svg">
	</a>
	<a href="https://opensource.org/licenses/Apache-2.0">
		<img alt="Apache-2.0 License" src="http://img.shields.io/badge/license-Apache-brightgreen.svg">
	</a>
</p>

**Worker** is a simple HTTP server that accepts FQL queries, executes them and returns their results.
OpenAPI v2 schema can be found [here](https://raw.githubusercontent.com/MontFerret/cli/master/reference/ferret-worker.yaml).

## Quick start

The Worker is shipped with dedicated Docker image that contains headless Google Chrome, so feel free to run queries using `cdp` driver:

DockerHub
```sh
docker run -d -p 8080:8080 montferret/worker
```
GitHub
```sh
docker run -d -p 8080:8080 ghcr.io/montferret/worker
```

Alternatively, if you want to use your own version of Chrome, you can run the Worker locally.

By installing the binary:

```shell
curl https://raw.githubusercontent.com/MontFerret/worker/master/install.sh | sh
worker
```

Or by building locally:

```sh
make
```

And then just make a POST request:

![worker](https://raw.githubusercontent.com/MontFerret/worker/master/assets/postman.png)

## System Resource Requirements
- 2 CPU
- 2 Gb of RAM

## Usage

### Endpoints

#### POST /
Executes a given query. The payload must have the following shape:

```
Query {
    text: String!
    params: Map<string, any>
}
```

#### GET /info
Returns a worker information that contains details about Chrome, Ferret and itself. Has the following shape:

```
Info {
    ip: String!
    version: Version! {
        worker: String!
        chrome: ChromeVersion! {
            browser: String!
            protocol: String!
            v8: String!
            webkit: String!
        }
        ferret: String!
    }
}
```


#### GET /health
Health check endpoint (for Kubernetes, e.g.). Returns empty 200.

### Run commands

```bash
  -cache-size uint
        amount of cached queries. 0 means no caching (default 100)
  -chrome-ip string
        Google Chrome remote IP address (default "127.0.0.1")
  -chrome-port uint
        Google Chrome remote debugging port (default 9222)
  -log-level string
        log level (default "debug")
  -port uint
        port to listen (default 8080)
  -version
        show version
  -help
        show this list

```
