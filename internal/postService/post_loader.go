package postservice

import (
	"context"
	pbl "post-service/genproto/post_loader_service"
	"post-service/internal/structs"

	"go.uber.org/zap"
)

func (s *service) LoadPosts(ctx context.Context, params *pbl.LoadPostParam) (*pbl.LoadingStatus, error) {

	var posts = []structs.Post{}
	err := s.postLoaderRepo.BatchInsert(ctx,posts)
	if err != nil {
		if err == structs.ErrNotFound {
			return nil, err
		}
		s.logger.Error("error while inserting posts into database", zap.Error(err))
		return nil, structs.ErrInternal
	}
	return &pbl.LoadingStatus{}, nil
}


func (s *service) GetJobStatus(context.Context, *pbl.JobId) (*pbl.LoadingStatus, error) {

	

}
