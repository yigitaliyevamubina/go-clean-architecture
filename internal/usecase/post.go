package usecase

import (
	"context"
	"fmt"
	"fourth-exam/post-service-clean-arch/internal/entity"
)

// PostUseCase -.
type PostUseCase struct {
	repo PostRepo
}

// New -.
func New(p PostRepo) *PostUseCase {
	return &PostUseCase{repo: p}
}

// Create Post
func (p *PostUseCase) CreatePost(ctx context.Context, req *entity.Post) (*entity.Post, error) {
	post, err := p.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("PostUseCase - Create - p.repo: %w", err)
	}

	return post, nil
}

// Get Post
func (p *PostUseCase) GetPost(ctx context.Context, id string) (*entity.Post, error) {
	post, err := p.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("PostUseCase - Get - p.repo: %w", err)
	}

	return post, nil
}

// Update Post
func (p *PostUseCase) UpdatePost(ctx context.Context, req *entity.Post) (*entity.Post, error) {
	post, err := p.repo.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("PostUseCase - Update - p.repo: %w", err)
	}

	return post, nil
}

// Delete Post
func (p *PostUseCase) DeletePost(ctx context.Context, id string) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("PostUseCase - Delete - p.repo: %w", err)
	}

	return nil
}

// List posts
func (p *PostUseCase) ListPosts(ctx context.Context, req *entity.GetListFilter) (*entity.Posts, error) {
	posts, err := p.repo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("PostUseCase - List - p.repo: %w", err)
	}

	return posts, nil
}
