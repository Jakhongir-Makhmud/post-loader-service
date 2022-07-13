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

// Warning max number of pages which has in requesting service not taken into account
func (s *service) LoadPosts(ctx context.Context, params *pbl.LoadPostParam) (*pbl.LoadingStatus, error) {
	var processingPages, lastProcessedPages int64
	// default loading pages is 50
	if params.Pages == 0 {
		params.Pages = 50
	}
	// here is used in memory cache which is thread safe, so we can use it concurrency without much worring
	if v, ok := s.cache.Get(pagesDoneKey); ok {
		lastProcessedPages = v.(int64)
		processingPages = lastProcessedPages + params.Pages
		s.cache.Set(pagesDoneKey, v.(int64)+params.Pages)
	} else {
		processingPages = params.Pages
		s.cache.Set(pagesDoneKey, params.Pages)
	}

	proccessId := fmt.Sprintf("%d", time.Now().Nanosecond())

	s.cache.Set(proccessId, &pbl.LoadingStatus{ProcessId: proccessId, WorkOfDone: 0, Status: "on process", TotalWork: params.Pages})

	jobs := make([]func(), 0, params.Pages)
	// ex: lastProcessedPages was 50, so we load posts till 50, next time we want load next 30 pages from 50 till 80,
	// we add last and current and got 80 (50 + 30).
	s.logger.Debug("pages", zap.Int64("lastpages:", lastProcessedPages), zap.Int64("current pages", processingPages))
	for i := lastProcessedPages; i < processingPages; i++ {
		job := func() {

			defer func() {
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
			posts, err := s.postSource.GetPostPage(int(i) + 1)
			if err != nil {
				s.logger.Error("error while getting post", zap.Error(err))
				return
			}

			err = s.postLoaderRepo.BatchInsert(ctx, posts)
			if err != nil {
				s.logger.Error("can't insert posts into database", zap.Error(err))
			}

		}
		jobs = append(jobs, job)
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
