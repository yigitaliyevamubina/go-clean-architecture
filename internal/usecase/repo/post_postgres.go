package repo

import (
	"context"
	"database/sql"
	"fmt"
	"fourth-exam/post-service-clean-arch/internal/entity"
	"fourth-exam/post-service-clean-arch/pkg/postgres"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// PostRepo -.
type PostRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *PostRepo {
	return &PostRepo{pg}
}

// Create post -.
func (p *PostRepo) Create(ctx context.Context, req *entity.Post) (*entity.Post, error) {
	if req.Id == "" {
		req.Id = uuid.New().String()
	}

	query, args, err := p.Builder.Insert("posts").
		Columns(`
			id,
			user_id,
			content,
			title,
			likes,
			dislikes,
			views,
			category,
			created_at
		`).
		Values(
			req.Id, req.UserId, req.Content, req.Title,
			req.Likes, req.Dislikes, req.Views, req.Category, time.Now()).Suffix(
		`RETURNING created_at, updated_at`,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostRepo - CreatePost - p.Builder: %w", err)
	}
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)


	row := p.Pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("PostRepo - CreatePost row.Scan: %v", err)
	}

	req.CreatedAt = createdAt.String()
	if updatedAt.Valid {
		req.UpdatedAt = updatedAt.Time.String()
	}

	return req, nil
}

// Get Post -.
func (p *PostRepo) Get(ctx context.Context, id string) (*entity.Post, error) {
	query := p.Builder.
		Select(`
		id,
		user_id,
		content,
		title,
		likes,
		dislikes,
		views,
		category,
		created_at,
		updated_at
		`).
		From("posts")

	if id != "" {
		query = query.Where(squirrel.Eq{"id": id})
	} else {
		return nil, fmt.Errorf("id is required")
	}

	q, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostRepo - GetPost - p.Builder: %w", err)
	}

	var (
		post      entity.Post
		createdAt time.Time
		updatedAt sql.NullTime
	)

	row := p.Pool.QueryRow(ctx, q, args...)
	if err := row.Scan(&post.Id, &post.UserId, &post.Content,
		&post.Title, &post.Likes, &post.Dislikes, &post.Views, &post.Category, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("PostRepo - GetPost row.Scan: %w", err)
	}

	post.CreatedAt = createdAt.String()
	if updatedAt.Valid {
		post.UpdatedAt = updatedAt.Time.String()
	}

	return &post, nil

}

// Update Post -.
func (p *PostRepo) Update(ctx context.Context, req *entity.Post) (*entity.Post, error) {
	var (
		updateMap = make(map[string]interface{})
		where     = squirrel.Eq{"id": req.Id}
	)

	updateMap["user_id"] = req.UserId
	updateMap["content"] = req.Content
	updateMap["title"] = req.Title
	updateMap["likes"] = req.Likes
	updateMap["dislikes"] = req.Dislikes
	updateMap["views"] = req.Views
	updateMap["category"] = req.Category
	updateMap["updated_at"] = time.Now()

	query := p.Builder.Update("posts").SetMap(updateMap).Where(where).Suffix("RETURNING created_at, updated_at")
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)

	q, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostRepo - UpdatePost - p.Builder: %w", err)
	}
	row := p.Pool.QueryRow(ctx, q, args...)
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("PostRepo - GetPost row.Scan: %w", err)
	}

	req.CreatedAt = createdAt.String()
	if updatedAt.Valid {
		req.UpdatedAt = updatedAt.Time.String()
	}

	return req, nil
}

// Delete Post -,
func (p *PostRepo) Delete(ctx context.Context, id string) error {
	query := p.Builder.Delete("posts")
	if id != "" {
		query = query.Where(squirrel.Eq{"id": id})
	} else {
		return fmt.Errorf("id is required")
	}

	q, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("PostRepo - DeletePost - p.Builder: %w", err)
	}

	_, err = p.Pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("PostRepo - DeletePost row Exec: %w", err)
	}

	return nil
}

// List Posts -.
func (p *PostRepo) List(ctx context.Context, req *entity.GetListFilter) (*entity.Posts, error) {
	query := p.Builder.
		Select(`
		id,
		user_id,
		content,
		title,
		likes,
		dislikes,
		views,
		category,
		created_at,
		updated_at
		`).
		From("posts")

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))
	query = query.OrderBy(req.OrderBy)

	if req.UserId != "" {
		query = query.Where(squirrel.Eq{"user_id": req.UserId})
	}

	q, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostRepo - ListPost - p.Builder: %w", err)
	}

	rows, err := p.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("PostRepo - PostHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var (
		posts = entity.Posts{Count: 0}
	)

	for rows.Next() {
		var (
			createdAt time.Time
			updatedAt sql.NullTime
			post      entity.Post
		)
		if err := rows.Scan(&post.Id, &post.UserId, &post.Content,
			&post.Title, &post.Likes, &post.Dislikes, &post.Views, &post.Category, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("PostRepo - ListPost row.Scan: %w", err)
		}

		post.CreatedAt = createdAt.String()
		if updatedAt.Valid {
			post.UpdatedAt = updatedAt.Time.String()
		}

		posts.Count++
		posts.Items = append(posts.Items, &post)
	}

	return &posts, nil
}
