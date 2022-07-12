package post_loader_repo

import (
	"context"
	"post-loader-service/internal/structs"
	"post-loader-service/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type PostLoaderRepo interface {
	BatchInsert(ctx context.Context, posts []structs.Post) error
}

type repo struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPosLoadertRepo(db *sqlx.DB, logger logger.Logger) PostLoaderRepo {
	return &repo{
		db:     db,
		logger: logger,
	}
}
