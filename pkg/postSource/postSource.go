package postSource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"post-loader-service/internal/structs"
	"post-loader-service/pkg/config"
	"post-loader-service/pkg/logger"
	"strconv"

	"go.uber.org/zap"
)

type PostSource interface {
	GetPostPage(page int) ([]structs.Post, error)
}

type source struct {
	url    string
	logger logger.Logger
}

func NewPostSource(cfg config.Config, l logger.Logger) PostSource {
	return &source{url: cfg.GetString("app.postSource.url"), logger: l}
}

func (s *source) GetPostPage(page int) ([]structs.Post, error) {
	pageStr := strconv.Itoa(page)
	client := &http.Client{}
	req, err := http.NewRequest("GET", s.url+"?page="+pageStr, nil)
	if err != nil {
		s.logger.Error("error while creating new request", zap.Error(err))
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("failed when requesting for posts", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	var posts DraftPost
	err = json.NewDecoder(resp.Body).Decode(&posts)

	if err != nil {
		s.logger.Error("can't unmarsha json to struct", zap.Error(err))
		b, err := ioutil.ReadAll(resp.Body)
		fmt.Println("\n\nbody in string", string(b))
		return nil, err
	}

	return s.transform(posts), nil
}

func (s *source) transform(draftPosts DraftPost) []structs.Post {

	newPosts := make([]structs.Post, len(draftPosts.Data))

	for i, p := range draftPosts.Data {
		post := structs.Post{Id: p.ID, Title: p.Title, Body: p.Body}

		newPosts[i] = post
	}
	return newPosts
}
