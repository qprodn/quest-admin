package organization

import (
	"context"

	v1 "quest-admin/api/gen/organization/v1"
	biz "quest-admin/internal/biz/organization"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostService struct {
	v1.UnimplementedPostServiceServer
	pc  *biz.PostUsecase
	log *log.Helper
}

func NewPostService(pc *biz.PostUsecase, logger log.Logger) *PostService {
	return &PostService{
		pc:  pc,
		log: log.NewHelper(logger),
	}
}

func (s *PostService) CreatePost(ctx context.Context, in *v1.CreatePostRequest) (*emptypb.Empty, error) {
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

func (s *PostService) GetPost(ctx context.Context, in *v1.GetPostRequest) (*v1.GetPostReply, error) {
	post, err := s.pc.GetPost(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetPostReply{
		Post: s.toProtoPost(post),
	}, nil
}

func (s *PostService) ListPosts(ctx context.Context, in *v1.ListPostsRequest) (*v1.ListPostsReply, error) {
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

func (s *PostService) UpdatePost(ctx context.Context, in *v1.UpdatePostRequest) (*emptypb.Empty, error) {
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

func (s *PostService) DeletePost(ctx context.Context, in *v1.DeletePostRequest) (*emptypb.Empty, error) {
	err := s.pc.DeletePost(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PostService) toProtoPost(post *biz.Post) *v1.PostInfo {
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
