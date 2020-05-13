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
   <a href="https://microbadger.com/images/montferret/worker">
      <img alt="Dockerimages" src="https://images.microbadger.com/badges/version/montferret/worker.svg">
   </a>
	<a href="http://opensource.org/licenses/MIT">
		<img alt="MIT License" src="http://img.shields.io/badge/license-MIT-brightgreen.svg">
	</a>
</p>

**Worker** is a simple HTTP server that accepts FQL queries, executes them and returns their results.

## Quick start

The Worker is shipped with dedicated Docker image that contains headless Google Chrome, so feel free to run queries using `cdp` driver:

```.env
docker run -p 8080:8080 -it montferret/worker
```

Alternatively, if you want to use your own version of Chrome, you can run the Worker locally:

```sh
make
```

And then just make a POST request:

![worker](https://raw.githubusercontent.com/MontFerret/worker/master/assets/postman.png)