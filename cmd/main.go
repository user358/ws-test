package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"ws-test/internal/balance"
	"ws-test/internal/grpc"
	"ws-test/internal/http"
	"ws-test/internal/log"
	"ws-test/internal/notification"
)

// application config
type config struct {
	HTTPPort int `envconfig:"HTTP_PORT" default:"8080"`
	GRPCPort int `envconfig:"GRPC_PORT" default:"8081"`
}

func main() {
	// context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initialize logger
	defer log.Init()()

	cfg := &config{}
	if err := envconfig.Process("WS_TEST", cfg); err != nil {
		log.Fatalw("failed to load config", "err", err)
	}

	// dependencies initialization
	db := balance.NewInMemoryStore()
	balanceRepository := balance.NewRepository(db)
	balanceService := balance.NewService(balanceRepository)
	notificationService := notification.NewService()

	wg := sync.WaitGroup{}

	// starting http server
	httpServer := http.New(cfg.HTTPPort, balanceService, notificationService)
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Run(ctx)
	}()

	// starting grpc server
	grpcServer := grpc.New(cfg.GRPCPort, balanceService)
	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Run(ctx)
	}()

	// wait until both servers stops
	wg.Wait()
}
