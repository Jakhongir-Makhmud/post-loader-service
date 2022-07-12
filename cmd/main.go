package main

import (
	"net"
	pbl "post-loader-service/genproto/post_loader_service" // protobuffer post loader service
	postLoaderService "post-loader-service/internal/postLoaderService"
	"post-loader-service/pkg/cache"
	"post-loader-service/pkg/config"
	"post-loader-service/pkg/db"
	"post-loader-service/pkg/logger"
	"post-loader-service/pkg/postSource"
	"post-loader-service/pkg/workerPool"
	postLoaderRepo "post-loader-service/repo"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg := config.NewConfig()

	db, logger := db.NewDB(cfg), logger.New(cfg.GetString("app.log.level"), cfg.GetString("app.name"))
	cache := cache.NewCache(cfg)

	postLoaderRepo := postLoaderRepo.NewPosLoadertRepo(db, logger)

	postSource := postSource.NewPostSource(cfg, logger)

	pool := workerPool.NewWorkerPool(runtime.GOMAXPROCS(0)*3, cfg.GetInt("app.workerPool.queueSize"))

	params := postLoaderService.Params{
		Logger:         logger,
		PostLoaderRepo: postLoaderRepo,
		Cache:          cache,
		WorkerPool:     pool,
		PostSource:     postSource,
	}
	service := postLoaderService.NewPostLoaderService(params)
	go pool.Run()
	listener, err := net.Listen("tcp", cfg.GetString("app.port"))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pbl.RegisterPostLoaderServiceServer(s, service)
	reflection.Register(s)

	logger.Info("service has stated it's job")

	panic(s.Serve(listener))
}
