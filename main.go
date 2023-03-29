package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-waitfor/waitfor"
	http "github.com/go-waitfor/waitfor-http"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/MontFerret/worker/internal/controllers"
	"github.com/MontFerret/worker/internal/server"
	"github.com/MontFerret/worker/internal/storage"
	"github.com/MontFerret/worker/pkg/caching"
	"github.com/MontFerret/worker/pkg/worker"
)

var (
	version string

	ferretVersion string

	port = flag.Uint64("port", 8080, "port to listen")

	noChrome = flag.Bool("no-chrome", false, "disable Chrome driver")

	chromeIP = flag.String("chrome-ip", "127.0.0.1", "Google Chrome remote IP address")

	chromeDebuggingPort = flag.Uint64("chrome-port", 9222, "Google Chrome remote debugging port")

	logLevel = flag.String(
		"log-level",
		zerolog.DebugLevel.String(),
		"log level",
	)

	cacheSize = flag.Uint(
		"cache-size",
		100,
		"amount of cached queries. 0 means no caching.",
	)

	requestLimit = flag.Uint64(
		"request-limit",
		0,
		"amount of requests per second for each IP. 0 means no limit.",
	)

	requestLimitTimeWindow = flag.Uint64(
		"request-limit-time-window",
		60*3,
		"amount of seconds for request rate limit time window",
	)

	bodyLimit = flag.Uint64(
		"body-limit",
		1000,
		"maximum size of request body in kb. 0 means no limit.",
	)

	showVersion = flag.Bool(
		"version",
		false,
		"show version",
	)

	help = flag.Bool(
		"help",
		false,
		"show this list",
	)
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("Worker: %s\n", version)
		fmt.Printf("Ferret: %s\n", ferretVersion)
		os.Exit(0)
	}

	z, err := zerolog.ParseLevel(*logLevel)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	level, _ := lecho.MatchZeroLevel(z)
	logger := lecho.New(
		os.Stdout,
		lecho.WithTimestamp(),
		lecho.WithLevel(level),
	)

	cdp := worker.CDPSettings{
		Host:     *chromeIP,
		Port:     *chromeDebuggingPort,
		Disabled: *noChrome,
	}

	if err := waitForChrome(cdp); err != nil {
		logger.Fatalf("wait for Chrome: %s", err)
	}

	srv, err := server.New(logger, server.Options{
		RequestLimit:           *requestLimit,
		RequestLimitTimeWindow: *requestLimitTimeWindow,
		BodyLimit:              *bodyLimit,
	})

	if err != nil {
		logger.Fatal(err)
	}

	cache, err := storage.NewCache(caching.WithSize(*cacheSize))

	if err != nil {
		logger.Fatal(errors.Wrap(err, "create cache storage"))
	}

	opts := make([]worker.Option, 0, 2)
	opts = append(opts, worker.WithCache(cache))

	if !cdp.Disabled {
		opts = append(opts, worker.WithCustomCDP(cdp))
	}

	wkr, err := worker.New(opts...)

	if err != nil {
		logger.Fatal(errors.Wrap(err, "create a worker instance"))
	}

	if err := setupControllers(srv, cdp, wkr); err != nil {
		logger.Fatal(err)
	}

	if err := srv.Run(*port); err != nil {
		logger.Fatal(err)
	}
}

func waitForChrome(cdp worker.CDPSettings) error {
	if cdp.Disabled {
		return nil
	}

	runner := waitfor.New(http.Use())

	return runner.Test(context.Background(), []string{
		cdp.BaseURL(),
	}, waitfor.WithAttempts(10))
}

func setupControllers(server *server.Server, cdp worker.CDPSettings, worker *worker.Worker) error {
	workerCtl, err := controllers.NewWorker(worker)

	if err != nil {
		return err
	}

	workerCtl.Use(server.Router())

	healthCtl, err := controllers.NewHealth(cdp)

	if err != nil {
		return err
	}

	healthCtl.Use(server.Router())

	infoCtl, err := controllers.NewInfo(controllers.InfoSettings{
		Version:       version,
		FerretVersion: ferretVersion,
		CDP:           cdp,
	})

	if err != nil {
		return err
	}

	infoCtl.Use(server.Router())

	return nil
}
