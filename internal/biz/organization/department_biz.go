package organization

import (
	"context"
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
)

type DepartmentRepo interface {
	Create(ctx context.Context, dept *Department) (*Department, error)
	FindByID(ctx context.Context, id string) (*Department, error)
	FindByName(ctx context.Context, name string) (*Department, error)
	List(ctx context.Context) ([]*Department, error)
	FindByParentID(ctx context.Context, parentID string) ([]*Department, error)
	Update(ctx context.Context, dept *Department) (*Department, error)
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

	existing, err := uc.repo.FindByName(ctx, dept.Name)
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
	return uc.repo.FindByID(ctx, id)
}

func (uc *DepartmentUsecase) GetDepartmentTree(ctx context.Context) ([]*Department, error) {
	uc.log.WithContext(ctx).Infof("GetDepartmentTree")

	departments, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	deptMap := make(map[string]*Department)
	for _, dept := range departments {
		deptMap[dept.ID] = dept
	}

	var roots []*Department
	for _, dept := range deptMap {
		if dept.ParentID == "" || dept.ParentID == "0" {
			roots = append(roots, dept)
		} else if parent, ok := deptMap[dept.ParentID]; ok {
			parent.Children = append(parent.Children, dept)
		}
	}

	return roots, nil
}

func (uc *DepartmentUsecase) UpdateDepartment(ctx context.Context, dept *Department) (*Department, error) {
	uc.log.WithContext(ctx).Infof("UpdateDepartment: id=%s, name=%s", dept.ID, dept.Name)

	_, err := uc.repo.FindByID(ctx, dept.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, dept)
}

func (uc *DepartmentUsecase) DeleteDepartment(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteDepartment: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	children, err := uc.repo.FindByParentID(ctx, id)
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
