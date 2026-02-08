package organization

import (
	"context"
	"quest-admin/internal/data/idgen"
	"quest-admin/pkg/errorx"
	"quest-admin/pkg/util/pagination"
	"quest-admin/types/consts/id"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/log"
)

type PostRepo interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	FindByName(ctx context.Context, name string) (*Post, error)
	FindByCode(ctx context.Context, code string) (*Post, error)
	List(ctx context.Context, opt *WherePostOpt) ([]*Post, error)
	Count(ctx context.Context, opt *WherePostOpt) (int64, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
	FindListByIDs(ctx context.Context, ids []string) ([]*Post, error)
}

type PostUsecase struct {
	idgen *idgen.IDGenerator
	repo  PostRepo
	log   *log.Helper
}

func NewPostUsecase(
	idgen *idgen.IDGenerator,
	repo PostRepo,
	logger log.Logger,
) *PostUsecase {
	return &PostUsecase{
		idgen: idgen,
		repo:  repo,
		log:   log.NewHelper(log.With(logger, "module", "organization/biz/post")),
	}
}

func (uc *PostUsecase) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	uc.log.WithContext(ctx).Infof("CreatePost: name=%s, code=%s", post.Name, post.Code)

	existing, err := uc.repo.FindByName(ctx, post.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errorx.Err(errkey.ErrPostNameExists)
	}
	post.ID = uc.idgen.NextID(id.POST)

	return uc.repo.Create(ctx, post)
}

func (uc *PostUsecase) GetPost(ctx context.Context, id string) (*Post, error) {
	uc.log.WithContext(ctx).Infof("GetPost: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *PostUsecase) ListPosts(ctx context.Context, query *ListPostsQuery) (*ListPostsResult, error) {
	uc.log.WithContext(ctx).Infof("ListPosts: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)

	opt := &WherePostOpt{
		Limit:     query.PageSize,
		Offset:    pagination.GetOffset(query.Page, query.PageSize),
		Keyword:   query.Keyword,
		Status:    query.Status,
		SortField: query.SortField,
		SortOrder: query.SortOrder,
	}

	list, err := uc.repo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询岗位列表失败", err)
		return nil, err
	}

	total, err := uc.repo.Count(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询岗位列表总数失败", err)
		return nil, err
	}

	return &ListPostsResult{
		Posts:      list,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: pagination.GetTotalPages(total, int64(query.PageSize)),
	}, nil
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
		return errorx.Err(errkey.ErrPostHasUsers)
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *PostUsecase) ListByPostIDs(ctx context.Context, postIds []string) ([]*Post, error) {
	var res []*Post
	if len(postIds) == 0 {
		return res, nil
	}
	posts, err := uc.repo.FindListByIDs(ctx, postIds)
	if err != nil {
		return posts, err
	}
	return posts, nil
}
