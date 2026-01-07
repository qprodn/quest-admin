package organization

import (
	"context"
	"time"

	v1 "quest-admin/api/gen/organization/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrDepartmentNotFound      = errors.NotFound(v1.ErrorReason_DEPARTMENT_NOT_FOUND.String(), "department not found")
	ErrDepartmentNameExists    = errors.Conflict(v1.ErrorReason_DEPARTMENT_NAME_EXISTS.String(), "department name already exists")
	ErrDepartmentHasChildren   = errors.BadRequest(v1.ErrorReason_DEPARTMENT_HAS_CHILDREN.String(), "department has children")
	ErrDepartmentHasUsers      = errors.BadRequest(v1.ErrorReason_DEPARTMENT_HAS_USERS.String(), "department has users")
	ErrInvalidParentDepartment = errors.BadRequest(v1.ErrorReason_INVALID_PARENT_DEPARTMENT.String(), "invalid parent department")
	ErrPostNotFound            = errors.NotFound(v1.ErrorReason_POST_NOT_FOUND.String(), "post not found")
	ErrPostNameExists          = errors.Conflict(v1.ErrorReason_POST_NAME_EXISTS.String(), "post name already exists")
	ErrPostHasUsers            = errors.BadRequest(v1.ErrorReason_POST_HAS_USERS.String(), "post has users")
)

type Department struct {
	ID           string
	Name         string
	ParentID     string
	Sort         int32
	LeaderUserID string
	Phone        string
	Email        string
	Status       int32
	CreateBy     string
	CreateAt     time.Time
	UpdateBy     string
	UpdateAt     time.Time
	TenantID     string
	Children     []*Department
}

type Post struct {
	ID       string
	Name     string
	Code     string
	Sort     int32
	Status   int32
	Remark   string
	CreateBy string
	CreateAt time.Time
	UpdateBy string
	UpdateAt time.Time
	TenantID string
}

type ListPostsQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type ListPostsResult struct {
	Posts      []*Post
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type DepartmentRepo interface {
	Create(ctx context.Context, dept *Department) (*Department, error)
	GetByID(ctx context.Context, id string) (*Department, error)
	GetByName(ctx context.Context, name string) (*Department, error)
	GetTree(ctx context.Context) ([]*Department, error)
	GetChildren(ctx context.Context, parentID string) ([]*Department, error)
	Update(ctx context.Context, dept *Department) (*Department, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
}

type PostRepo interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	GetByID(ctx context.Context, id string) (*Post, error)
	GetByName(ctx context.Context, name string) (*Post, error)
	GetByCode(ctx context.Context, code string) (*Post, error)
	List(ctx context.Context, query *ListPostsQuery) (*ListPostsResult, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
}

type DepartmentUsecase struct {
	repo DepartmentRepo
	log  *log.Helper
}

func NewDepartmentUsecase(repo DepartmentRepo, logger log.Logger) *DepartmentUsecase {
	return &DepartmentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *DepartmentUsecase) CreateDepartment(ctx context.Context, dept *Department) (*Department, error) {
	uc.log.WithContext(ctx).Infof("CreateDepartment: name=%s, parentID=%s", dept.Name, dept.ParentID)

	existing, err := uc.repo.GetByName(ctx, dept.Name)
	if err != nil && !errors.Is(err, ErrDepartmentNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDepartmentNameExists
	}

	return uc.repo.Create(ctx, dept)
}

func (uc *DepartmentUsecase) GetDepartment(ctx context.Context, id string) (*Department, error) {
	uc.log.WithContext(ctx).Infof("GetDepartment: id=%s", id)
	return uc.repo.GetByID(ctx, id)
}

func (uc *DepartmentUsecase) GetDepartmentTree(ctx context.Context) ([]*Department, error) {
	uc.log.WithContext(ctx).Infof("GetDepartmentTree")
	return uc.repo.GetTree(ctx)
}

func (uc *DepartmentUsecase) UpdateDepartment(ctx context.Context, dept *Department) (*Department, error) {
	uc.log.WithContext(ctx).Infof("UpdateDepartment: id=%s, name=%s", dept.ID, dept.Name)

	_, err := uc.repo.GetByID(ctx, dept.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, dept)
}

func (uc *DepartmentUsecase) DeleteDepartment(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteDepartment: id=%s", id)

	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	children, err := uc.repo.GetChildren(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrDepartmentHasChildren
	}

	hasUsers, err := uc.repo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if hasUsers {
		return ErrDepartmentHasUsers
	}

	return uc.repo.Delete(ctx, id)
}

type PostUsecase struct {
	repo PostRepo
	log  *log.Helper
}

func NewPostUsecase(repo PostRepo, logger log.Logger) *PostUsecase {
	return &PostUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *PostUsecase) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	uc.log.WithContext(ctx).Infof("CreatePost: name=%s, code=%s", post.Name, post.Code)

	existing, err := uc.repo.GetByName(ctx, post.Name)
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
	return uc.repo.GetByID(ctx, id)
}

func (uc *PostUsecase) ListPosts(ctx context.Context, query *ListPostsQuery) (*ListPostsResult, error) {
	uc.log.WithContext(ctx).Infof("ListPosts: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)
	return uc.repo.List(ctx, query)
}

func (uc *PostUsecase) UpdatePost(ctx context.Context, post *Post) (*Post, error) {
	uc.log.WithContext(ctx).Infof("UpdatePost: id=%s, name=%s", post.ID, post.Name)

	_, err := uc.repo.GetByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, post)
}

func (uc *PostUsecase) DeletePost(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeletePost: id=%s", id)

	_, err := uc.repo.GetByID(ctx, id)
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
