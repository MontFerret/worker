package main

import (
	"context"
	"flag"

	"github.com/MontFerret/worker/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/waitfor/pkg/runner"
	waitrunner "github.com/ziflex/waitfor/pkg/runner"
)

var (
	port               = flag.String("port", "8080", "port to listen")
	chromeDebugginPort = flag.String("chrome-port", "9222", "Google Chrome remote debugging port")
)

func main() {
	flag.Parse()

	err := waitForChrome(*chromeDebugginPort)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("wait for Chrome")
	}

	server := server.New()

	log.Info().Msgf("listen at :%s", *port)

	err = server.Run(*port)
	log.Err(err).
		Timestamp().
		Msg("listen and server")
}

func waitForChrome(port string) error {
	return waitrunner.Test(context.Background(), []string{
		"http://127.0.0.1:" + port,
	}, runner.WithAttempts(10))
}
