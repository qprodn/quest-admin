package organization

import (
	"context"

	v1 "quest-admin/api/gen/organization/v1"
	biz "quest-admin/internal/biz/organization"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrganizationService struct {
	v1.UnimplementedOrganizationServiceServer
	dc  *biz.DepartmentUsecase
	pc  *biz.PostUsecase
	log *log.Helper
}

func NewOrganizationService(dc *biz.DepartmentUsecase, pc *biz.PostUsecase, logger log.Logger) *OrganizationService {
	return &OrganizationService{
		dc:  dc,
		pc:  pc,
		log: log.NewHelper(logger),
	}
}

func (s *OrganizationService) CreateDepartment(ctx context.Context, in *v1.CreateDepartmentRequest) (*emptypb.Empty, error) {
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

func (s *OrganizationService) GetDepartment(ctx context.Context, in *v1.GetDepartmentRequest) (*v1.GetDepartmentReply, error) {
	dept, err := s.dc.GetDepartment(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetDepartmentReply{
		Department: s.toProtoDepartment(dept),
	}, nil
}

func (s *OrganizationService) GetDepartmentTree(ctx context.Context, in *emptypb.Empty) (*v1.GetDepartmentTreeReply, error) {
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

func (s *OrganizationService) UpdateDepartment(ctx context.Context, in *v1.UpdateDepartmentRequest) (*emptypb.Empty, error) {
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

func (s *OrganizationService) DeleteDepartment(ctx context.Context, in *v1.DeleteDepartmentRequest) (*emptypb.Empty, error) {
	err := s.dc.DeleteDepartment(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *OrganizationService) CreatePost(ctx context.Context, in *v1.CreatePostRequest) (*emptypb.Empty, error) {
	post := &biz.Post{
		Name:   in.GetName(),
		Code:   in.GetCode(),
		Sort:   in.GetSort(),
		Remark: in.GetRemark(),
		Status: 1,
	}

	_, err := s.pc.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *OrganizationService) GetPost(ctx context.Context, in *v1.GetPostRequest) (*v1.GetPostReply, error) {
	post, err := s.pc.GetPost(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetPostReply{
		Post: s.toProtoPost(post),
	}, nil
}

func (s *OrganizationService) ListPosts(ctx context.Context, in *v1.ListPostsRequest) (*v1.ListPostsReply, error) {
	query := &biz.ListPostsQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Keyword:   in.GetKeyword(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.pc.ListPosts(ctx, query)
	if err != nil {
		return nil, err
	}

	posts := make([]*v1.PostInfo, 0, len(result.Posts))
	for _, post := range result.Posts {
		posts = append(posts, s.toProtoPost(post))
	}

	return &v1.ListPostsReply{
		Posts:      posts,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *OrganizationService) UpdatePost(ctx context.Context, in *v1.UpdatePostRequest) (*emptypb.Empty, error) {
	post := &biz.Post{
		ID:     in.GetId(),
		Name:   in.GetName(),
		Code:   in.GetCode(),
		Sort:   in.GetSort(),
		Status: in.GetStatus(),
		Remark: in.GetRemark(),
	}

	_, err := s.pc.UpdatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *OrganizationService) DeletePost(ctx context.Context, in *v1.DeletePostRequest) (*emptypb.Empty, error) {
	err := s.pc.DeletePost(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *OrganizationService) toProtoDepartment(dept *biz.Department) *v1.DepartmentInfo {
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

func (s *OrganizationService) toProtoPost(post *biz.Post) *v1.PostInfo {
	return &v1.PostInfo{
		Id:       post.ID,
		Name:     post.Name,
		Code:     post.Code,
		Sort:     post.Sort,
		Status:   post.Status,
		Remark:   post.Remark,
		CreateAt: timestamppb.New(post.CreateAt),
		UpdateAt: timestamppb.New(post.UpdateAt),
		TenantId: post.TenantID,
	}
}
