package postLoaderservice

import (
	"post-loader-service/pkg/cache"
	"post-loader-service/pkg/logger"
	"post-loader-service/pkg/postSource"
	"post-loader-service/pkg/workerPool"
	post_repo "post-loader-service/repo"
)

type Params struct {
	PostLoaderRepo post_repo.PostLoaderRepo
	PostSource     postSource.PostSource
	WorkerPool     workerPool.WorkerPoll
	Cache          cache.Cache
	Logger         logger.Logger
}

type service struct {
	postLoaderRepo post_repo.PostLoaderRepo
	postSource     postSource.PostSource
	workerPool     workerPool.WorkerPoll
	cache          cache.Cache
	logger         logger.Logger
}

func NewPostLoaderService(p Params) *service {
	return &service{
		postLoaderRepo: p.PostLoaderRepo,
		postSource:     p.PostSource,
		workerPool:     p.WorkerPool,
		logger:         p.Logger,
		cache:          p.Cache,
	}
}
