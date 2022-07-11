package postservice

import (
	"post-service/pkg/logger"
	post_repo "post-service/repo"
)

type service struct {
	postLoaderRepo post_repo.PostLoaderRepo
	logger   logger.Logger
}

func NewPostLoaderService(repo post_repo.PostLoaderRepo, logger logger.Logger) *service {
	return &service{
		postLoaderRepo: repo,
		logger:   logger,
	}
}
