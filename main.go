package main

import (
	"context"
	"fmt"

	"github.com/seb7887/heimdallr/config"
	"github.com/seb7887/heimdallr/health"
	"github.com/seb7887/heimdallr/server"
	"github.com/seb7887/heimdallr/server/grpc"
	"github.com/seb7887/heimdallr/service"
	"github.com/seb7887/heimdallr/storage"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		port     = config.GetConfig().HealthPort
		grpcPort = config.GetConfig().GRPCPort
		httpAddr = fmt.Sprintf(":%d", port)
		grpcAddr = fmt.Sprintf(":%d", grpcPort)

		repo          = storage.InitializeRepository()
		healthService = health.NewService(repo)
		grpcService   = service.NewService(repo)
	)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		srv := grpc.New(grpcAddr, grpcService)
		log.Infof("gRPC server running at %s", grpcAddr)
		return srv.Serve(ctx)
	})

	g.Go(func() error {
		httpSrv := server.New(healthService, httpAddr)
		log.Infof("HTTP server running at %s", httpAddr)
		return httpSrv.Serve(ctx)
	})

	log.Fatal(g.Wait())
}
