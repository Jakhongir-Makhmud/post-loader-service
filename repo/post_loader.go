package post_repo

import (
	"context"
	"post-service/internal/structs"

	"go.uber.org/zap"
)

func (r *repo) BatchInsert(ctx context.Context, posts []structs.Post) error {

	if len(posts) == 0 {
		r.logger.Warn("empty slice is got to insert")
		return structs.ErrNoData
	}

	query := `INSERT INTO posts ( post_id, title, body) VALUES (:post_id, :title, :body)`

	result, err := r.db.NamedExecContext(ctx, query, posts)
	if err != nil {
		r.logger.Error("error while inserting posts into database", zap.Error(err))
		return err
	}

	affected, err := result.RowsAffected()
	if affected == 0 && err != nil {
		r.logger.Error("something went wrong", zap.Error(err))
		return err
	}

	return nil

}
