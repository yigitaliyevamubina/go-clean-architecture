package usecase_test

import (
	"context"
	"errors"
	"fourth-exam/post-service-clean-arch/internal/entity"
	"fourth-exam/post-service-clean-arch/internal/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	errInternalServerErr = errors.New("internal server error")
	errBadRequest = errors.New("bad request")
)


type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func post(t *testing.T) (*usecase.PostUseCase, *MockPostRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockPostRepo(mockCtl)

	post := usecase.New(repo)

	return post, repo
}

func TestPost(t *testing.T) {
	t.Parallel()

	post, repo := post(t)

	body := entity.Post{
		UserId: "",
		Content: "Content",
		Title: "Post title",
		Likes: 0,
		Dislikes: 0,
		Views: 10,
		Category: "Post category",
	}

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().List(context.Background(), &entity.GetListFilter{Page: 1, Limit: 10}).Return(nil, nil)
			},
			res: (*entity.Posts)(nil),
			err: nil,
		},
		{
			name: "result with error",
			mock: func() {
				repo.EXPECT().List(context.Background(), &entity.GetListFilter{Page: 1, Limit: 10}).Return(nil, errInternalServerErr)
			},
			res: (*entity.Posts)(nil),
			err: errInternalServerErr,
		},
	}

	createTest := test{
		name: "create with error",
		mock: func() {
			repo.EXPECT().Create(context.Background(), &body).Return(nil, errBadRequest)
		},
		res: (*entity.Post)(nil),
		err: errBadRequest,
	}

	for _, ts := range tests {
		ts := ts

		t.Run(ts.name, func(t *testing.T) {
			t.Parallel()

			ts.mock()

			res, err := post.ListPosts(context.Background(), &entity.GetListFilter{Page: 1, Limit: 10})

			require.Equal(t, ts.res, res)
			require.ErrorIs(t, err, ts.err)
		})
	}

	t.Run(createTest.name, func(t *testing.T) {
		t.Parallel()

		createTest.mock()

		res, err := post.CreatePost(context.Background(), &body)

		require.Equal(t, createTest.res, res)
		require.ErrorIs(t, err, createTest.err)
	})
}
