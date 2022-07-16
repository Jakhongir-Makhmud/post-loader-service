package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	pbl "post-loader-service/genproto/post_loader_service" // protobuffer post loader service
	postLoaderService "post-loader-service/internal/postLoaderService"
	"post-loader-service/pkg/cache"
	"post-loader-service/pkg/config"
	"post-loader-service/pkg/db"
	"post-loader-service/pkg/logger"
	"post-loader-service/pkg/postSource"
	"post-loader-service/pkg/workerPool"
	postLoaderRepo "post-loader-service/repo"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutDownCh := make(chan os.Signal, 1)
	signal.Notify(shutDownCh, os.Interrupt)

	cfg := config.NewConfig()

	dbConn, logger := db.NewDB(cfg), logger.New(cfg.GetString("app.log.level"), cfg.GetString("app.name"))
	defer dbConn.Close()
	cache := cache.NewCache(cfg)

	postLoaderRepo := postLoaderRepo.NewPosLoadertRepo(dbConn, logger)

	postSource := postSource.NewPostSource(cfg, logger)
	// max number of workers is small because of requesting server seems don't support cocurrent request, so by this we can load posts without any errors.
	pool := workerPool.NewWorkerPool(4, cfg.GetInt("app.workerPool.queueSize"))

	params := postLoaderService.Params{
		Logger:         logger,
		PostLoaderRepo: postLoaderRepo,
		Cache:          cache,
		WorkerPool:     pool,
		PostSource:     postSource,
	}
	service := postLoaderService.NewPostLoaderService(params)
	go pool.Run(ctx)

	listener, err := net.Listen("tcp", cfg.GetString("app.port"))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pbl.RegisterPostLoaderServiceServer(s, service)
	reflection.Register(s)

	go func() {
		<-shutDownCh // it will wait until signal is received
		s.GracefulStop() // stop server gracefully, stop accepting requests
		cancel()         // to stop workers
		logger.Info("service shuted down gracefully")
	}()
	logger.Info("service has started it's job on port: " + cfg.GetString("app.port"))

	err = s.Serve(listener)
	if err == nil {
		logger.Warn("server is stoped", zap.Error(err))
	} else {
		panic(err)
	}
}
