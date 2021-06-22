package main

import (
	"context"
	"fmt"

	"github.com/seb7887/heimdallr/config"
	"github.com/seb7887/heimdallr/health"
	"github.com/seb7887/heimdallr/server"
	"github.com/seb7887/heimdallr/storage"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		port     = config.GetConfig().HealthPort
		httpAddr = fmt.Sprintf(":%d", port)

		repo          = storage.InitializeRepository()
		healthService = health.NewService(repo)
	)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		httpSrv := server.New(healthService, httpAddr)
		log.Infof("HTTP server running at %s", httpAddr)
		return httpSrv.Serve(ctx)
	})

	log.Fatal(g.Wait())
}
