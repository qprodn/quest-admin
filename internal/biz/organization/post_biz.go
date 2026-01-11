package organization

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type PostRepo interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	FindByName(ctx context.Context, name string) (*Post, error)
	FindByCode(ctx context.Context, code string) (*Post, error)
	List(ctx context.Context, query *ListPostsQuery) (*ListPostsResult, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
}

type PostUsecase struct {
	repo PostRepo
	log  *log.Helper
}

func NewPostUsecase(repo PostRepo, logger log.Logger) *PostUsecase {
	return &PostUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "organization/biz/post")),
	}
}

func (uc *PostUsecase) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	uc.log.WithContext(ctx).Infof("CreatePost: name=%s, code=%s", post.Name, post.Code)

	existing, err := uc.repo.FindByName(ctx, post.Name)
	if err != nil && !errors.Is(err, ErrPostNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPostNameExists
	}

	return uc.repo.Create(ctx, post)
}

func (uc *PostUsecase) GetPost(ctx context.Context, id string) (*Post, error) {
	uc.log.WithContext(ctx).Infof("GetPost: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *PostUsecase) ListPosts(ctx context.Context, query *ListPostsQuery) (*ListPostsResult, error) {
	uc.log.WithContext(ctx).Infof("ListPosts: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)
	return uc.repo.List(ctx, query)
}

func (uc *PostUsecase) UpdatePost(ctx context.Context, post *Post) (*Post, error) {
	uc.log.WithContext(ctx).Infof("UpdatePost: id=%s, name=%s", post.ID, post.Name)

	_, err := uc.repo.FindByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, post)
}

func (uc *PostUsecase) DeletePost(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeletePost: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	hasUsers, err := uc.repo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if hasUsers {
		return ErrPostHasUsers
	}

	return uc.repo.Delete(ctx, id)
}
