package organization

import (
	"context"

	v1 "quest-admin/api/gen/organization/v1"
	biz "quest-admin/internal/biz/organization"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DepartmentService struct {
	v1.UnimplementedDepartmentServiceServer
	dc  *biz.DepartmentUsecase
	log *log.Helper
}

func NewDepartmentService(dc *biz.DepartmentUsecase, logger log.Logger) *DepartmentService {
	return &DepartmentService{
		dc:  dc,
		log: log.NewHelper(log.With(logger, "module", "organization/service")),
	}
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, in *v1.CreateDepartmentRequest) (*emptypb.Empty, error) {
	dept := &biz.Department{
		Name:         in.GetName(),
		ParentID:     in.GetParentId(),
		Sort:         in.GetSort(),
		LeaderUserID: in.GetLeaderUserId(),
		Phone:        in.GetPhone(),
		Email:        in.GetEmail(),
		Status:       1,
	}

	_, err := s.dc.CreateDepartment(ctx, dept)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DepartmentService) GetDepartment(ctx context.Context, in *v1.GetDepartmentRequest) (*v1.GetDepartmentReply, error) {
	dept, err := s.dc.GetDepartment(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetDepartmentReply{
		Department: s.toProtoDepartment(dept),
	}, nil
}

func (s *DepartmentService) GetDepartmentTree(ctx context.Context, in *emptypb.Empty) (*v1.GetDepartmentTreeReply, error) {
	departments, err := s.dc.GetDepartmentTree(ctx)
	if err != nil {
		return nil, err
	}

	protoDepartments := make([]*v1.DepartmentInfo, 0, len(departments))
	for _, dept := range departments {
		protoDepartments = append(protoDepartments, s.toProtoDepartment(dept))
	}

	return &v1.GetDepartmentTreeReply{
		Departments: protoDepartments,
	}, nil
}

func (s *DepartmentService) UpdateDepartment(ctx context.Context, in *v1.UpdateDepartmentRequest) (*emptypb.Empty, error) {
	dept := &biz.Department{
		ID:           in.GetId(),
		Name:         in.GetName(),
		ParentID:     in.GetParentId(),
		Sort:         in.GetSort(),
		LeaderUserID: in.GetLeaderUserId(),
		Phone:        in.GetPhone(),
		Email:        in.GetEmail(),
		Status:       in.GetStatus(),
	}

	_, err := s.dc.UpdateDepartment(ctx, dept)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DepartmentService) DeleteDepartment(ctx context.Context, in *v1.DeleteDepartmentRequest) (*emptypb.Empty, error) {
	err := s.dc.DeleteDepartment(ctx, in.GetId())
	if err != nil {
		s.log.WithContext(ctx).Errorw("failed to delete department", "id", in.GetId(), "error", err)
		return nil, err
	}

	s.log.WithContext(ctx).Infow("department deleted successfully", "id", in.GetId())
	return &emptypb.Empty{}, nil
}

func (s *DepartmentService) toProtoDepartment(dept *biz.Department) *v1.DepartmentInfo {
	children := make([]*v1.DepartmentInfo, 0, len(dept.Children))
	for _, child := range dept.Children {
		children = append(children, s.toProtoDepartment(child))
	}

	return &v1.DepartmentInfo{
		Id:           dept.ID,
		Name:         dept.Name,
		ParentId:     dept.ParentID,
		Sort:         dept.Sort,
		LeaderUserId: dept.LeaderUserID,
		Phone:        dept.Phone,
		Email:        dept.Email,
		Status:       dept.Status,
		CreateAt:     timestamppb.New(dept.CreateAt),
		UpdateAt:     timestamppb.New(dept.UpdateAt),
		TenantId:     dept.TenantID,
		Children:     children,
	}
}
