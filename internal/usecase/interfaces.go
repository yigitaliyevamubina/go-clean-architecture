package usecase

import (
	"context"
	"fourth-exam/post-service-clean-arch/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Post -. 
	Post interface {
		CreatePost(context.Context, *entity.Post) (*entity.Post, error)
		GetPost(context.Context, string) (*entity.Post, error) 
		UpdatePost(context.Context, *entity.Post) (*entity.Post, error)
		DeletePost(context.Context, string) (error) 
		ListPosts(context.Context, *entity.GetListFilter) (*entity.Posts, error)
	}

	// PostRepo -. 
	PostRepo interface {
		Create(context.Context, *entity.Post) (*entity.Post, error)
		Get(context.Context, string) (*entity.Post, error) 
		Update(context.Context, *entity.Post) (*entity.Post, error)
		Delete(context.Context, string) (error) 
		List(context.Context, *entity.GetListFilter) (*entity.Posts, error)
	}
)