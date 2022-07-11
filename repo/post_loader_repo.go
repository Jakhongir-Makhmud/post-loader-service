package post_repo

import (
	"context"
	"post-service/internal/structs"
	"post-service/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type PostLoaderRepo interface {
	BatchInsert(ctx context.Context, posts []structs.Post) error
}

type repo struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostRepo(db *sqlx.DB, logger logger.Logger) PostLoaderRepo {
	return &repo{
		db:     db,
		logger: logger,
	}
}
