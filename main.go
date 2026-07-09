package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-waitfor/waitfor"
	http "github.com/go-waitfor/waitfor-http"
	"github.com/labstack/echo/v4"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/MontFerret/worker/internal/controllers"
	"github.com/MontFerret/worker/internal/server"
	"github.com/MontFerret/worker/pkg/worker"
)

var (
	version string

	ferretVersion string

	defaultHTTPPolicy = worker.DefaultHTTPPolicy()

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

	bodyLimit = flag.String(
		"body-limit",
		"1M",
		"maximum allowed size for a request body (e.g., 4K, 4KB, 10M, 1G). Empty string means no limit.",
	)

	fsRoot = flag.String(
		"fs-root",
		"",
		"file system root directory for FQL IO::FS functions. Defaults to the current working directory.",
	)

	httpAllowedHosts = flag.String(
		"http-allowed-hosts",
		strings.Join(defaultHTTPPolicy.AllowedHosts, ","),
		"comma-separated exact hosts or host:port values allowed for Ferret HTTP requests",
	)

	httpAllowAllHosts = flag.Bool(
		"http-allow-all-hosts",
		defaultHTTPPolicy.AllowAllHosts,
		"allow Ferret HTTP requests to any host while still applying scheme, timeout, size, redirect, and literal-address policy",
	)

	httpBlockedHosts = flag.String(
		"http-blocked-hosts",
		strings.Join(defaultHTTPPolicy.BlockedHosts, ","),
		"comma-separated exact hosts or host:port values blocked for Ferret HTTP requests",
	)

	httpTimeout = flag.Duration(
		"http-timeout",
		defaultHTTPPolicy.Timeout,
		"timeout for Ferret HTTP requests",
	)

	httpMaxRequestSize = flag.Int64(
		"http-max-request-size",
		defaultHTTPPolicy.MaxRequestSize,
		"maximum Ferret HTTP request body size in bytes. 0 means no limit.",
	)

	httpMaxResponseSize = flag.Int64(
		"http-max-response-size",
		defaultHTTPPolicy.MaxResponseSize,
		"maximum Ferret HTTP response body size in bytes. 0 means no limit.",
	)

	httpMaxRedirects = flag.Int(
		"http-max-redirects",
		defaultHTTPPolicy.MaxRedirects,
		"maximum number of redirects followed by Ferret HTTP requests. 0 uses the Go standard library default.",
	)

	httpFollowRedirects = flag.Bool(
		"http-follow-redirects",
		defaultHTTPPolicy.FollowRedirects,
		"follow redirects for Ferret HTTP requests",
	)

	httpAllowLocalhost = flag.Bool(
		"http-allow-localhost",
		defaultHTTPPolicy.AllowLocalhost,
		"allow Ferret HTTP requests to localhost and loopback literal addresses",
	)

	httpAllowPrivateNetworks = flag.Bool(
		"http-allow-private-networks",
		defaultHTTPPolicy.AllowPrivateNetworks,
		"allow Ferret HTTP requests to private-network literal IP addresses",
	)

	httpBlockedRequestHeaders = flag.String(
		"http-blocked-request-headers",
		strings.Join(defaultHTTPPolicy.BlockedRequestHeaders, ","),
		"comma-separated request headers removed from Ferret HTTP requests",
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

type httpPolicyConfig struct {
	AllowedHosts          string
	BlockedHosts          string
	BlockedRequestHeaders string
	Timeout               time.Duration
	MaxRequestSize        int64
	MaxResponseSize       int64
	MaxRedirects          int
	FollowRedirects       bool
	AllowAllHosts         bool
	AllowLocalhost        bool
	AllowPrivateNetworks  bool
}

func currentHTTPPolicyConfig() httpPolicyConfig {
	return httpPolicyConfig{
		AllowedHosts:          *httpAllowedHosts,
		BlockedHosts:          *httpBlockedHosts,
		BlockedRequestHeaders: *httpBlockedRequestHeaders,
		Timeout:               *httpTimeout,
		MaxRequestSize:        *httpMaxRequestSize,
		MaxResponseSize:       *httpMaxResponseSize,
		MaxRedirects:          *httpMaxRedirects,
		FollowRedirects:       *httpFollowRedirects,
		AllowAllHosts:         *httpAllowAllHosts,
		AllowLocalhost:        *httpAllowLocalhost,
		AllowPrivateNetworks:  *httpAllowPrivateNetworks,
	}
}

func newHTTPPolicyFromConfig(config httpPolicyConfig) (worker.HTTPPolicy, error) {
	policy := worker.HTTPPolicy{
		AllowedSchemes:        worker.DefaultHTTPPolicy().AllowedSchemes,
		AllowedHosts:          splitCSV(config.AllowedHosts),
		BlockedHosts:          splitCSV(config.BlockedHosts),
		BlockedRequestHeaders: splitCSV(config.BlockedRequestHeaders),
		Timeout:               config.Timeout,
		MaxRequestSize:        config.MaxRequestSize,
		MaxResponseSize:       config.MaxResponseSize,
		MaxRedirects:          config.MaxRedirects,
		FollowRedirects:       config.FollowRedirects,
		AllowAllHosts:         config.AllowAllHosts,
		AllowLocalhost:        config.AllowLocalhost,
		AllowPrivateNetworks:  config.AllowPrivateNetworks,
	}

	if policy.AllowAllHosts && len(policy.AllowedHosts) > 0 {
		return policy, fmt.Errorf("http allowed hosts and allow-all hosts cannot both be set")
	}

	return policy, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		out = append(out, part)
	}

	return out
}

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

	resolvedFSRoot, err := resolveFSRoot(*fsRoot, flagIsSet("fs-root"))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	httpPolicy, err := newHTTPPolicyFromConfig(currentHTTPPolicyConfig())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		RequestLimitSkipper: func(c echo.Context) bool {
			return c.Path() == controllers.HealthPath
		},
		BodyLimit: *bodyLimit,
	})

	if err != nil {
		logger.Fatal(err)
	}

	opts := []worker.Option{
		worker.WithCacheSize(*cacheSize),
		worker.WithFSRoot(resolvedFSRoot),
		worker.WithHTTPPolicy(httpPolicy),
		worker.WithRESTModule(),
	}

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

func flagIsSet(name string) bool {
	isSet := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			isSet = true
		}
	})

	return isSet
}

func resolveFSRoot(root string, isSet bool) (string, error) {
	if !isSet {
		wd, err := os.Getwd()

		if err != nil {
			return "", errors.Wrap(err, "get current working directory")
		}

		return wd, nil
	}

	root = strings.TrimSpace(root)

	if root == "" {
		return "", errors.New("fs root cannot be empty")
	}

	return root, nil
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
