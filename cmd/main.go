package main

import (
	"net"
	pbps "post-service/genproto/post_service" // protobuffer post service
	postService "post-service/internal/postService"
	"post-service/pkg/config"
	"post-service/pkg/db"
	"post-service/pkg/logger"
	post_repo "post-service/repo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg := config.NewConfig()

	db, logger :=db.NewDB(cfg), logger.New(cfg.GetString("app.log.level"), cfg.GetString("app.name"))

	postRepo := post_repo.NewPostRepo(db, logger)

	service := postService.NewPostService(postRepo, logger)

	listener, err := net.Listen("tcp", cfg.GetString("app.port"))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pbps.RegisterPostServiceServer(s, service)
	reflection.Register(s)

	logger.Info("service has stated it's job")

	panic(s.Serve(listener))
}
