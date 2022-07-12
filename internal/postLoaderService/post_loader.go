package postLoaderservice

import (
	"context"
	"fmt"
	pbl "post-loader-service/genproto/post_loader_service"
	"post-loader-service/internal/structs"
	"time"

	"go.uber.org/zap"
)

const pagesDoneKey = "pagesDone"

func (s *service) LoadPosts(ctx context.Context, params *pbl.LoadPostParam) (*pbl.LoadingStatus, error) {

	if params.Pages == 0 {
		params.Pages = 50
	}

	if v, ok := s.cache.Get(pagesDoneKey); ok {
		s.cache.Set(pagesDoneKey, v.(int64) + params.Pages)
	} else {
		s.cache.Set(pagesDoneKey, params.Pages)
	}

	proccessId := fmt.Sprintf("%d",time.Now().Nanosecond())

	s.cache.Set(proccessId, &pbl.LoadingStatus{ProcessId: proccessId, WorkOfDone: 0, Status: "on process", TotalWork: params.Pages})

	jobs := make([]func(), params.Pages)
	for i := range jobs {
		jobs[i] = func()  {
			
				defer func ()  {
					if v, ok := s.cache.Get(proccessId); ok {
						status := v.(*pbl.LoadingStatus)
						status.WorkOfDone = status.WorkOfDone + 1
						if status.TotalWork == status.WorkOfDone {
							status.Status = "done"
						}
						s.cache.Set(proccessId, status)
					}
				}()

				ctx := context.Background()
				posts, err := s.postSource.GetPostPage(i+1)
				if err != nil {
					s.logger.Error("error while getting post", zap.Error(err))
					return
				}

				err = s.postLoaderRepo.BatchInsert(ctx, posts)
				if err != nil {
					s.logger.Error("can't insert posts into database")
				}

			}
	}


	for i := range jobs {
		s.workerPool.AddJob(jobs[i])
	}
	status, _ := s.cache.Get(proccessId)
	return status.(*pbl.LoadingStatus), nil
}

func (s *service) GetJobStatus(ctx context.Context, id *pbl.JobId) (*pbl.LoadingStatus, error) {
	status, ok := s.cache.Get(id.Id)
	if ok {
		v, ok := status.(*pbl.LoadingStatus)
		if ok {
			return v, nil
		}
		return nil, structs.ErrTypeCast
	}
	return nil, structs.ErrNotFound
}

