package user

import (
	"context"
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"

	v1 "quest-admin/api/gen/user/v1"
	biz "quest-admin/internal/biz/user"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) AssignUserPosts(ctx context.Context, in *v1.AssignUserPostRequest) (*emptypb.Empty, error) {
	bo := &biz.AssignUserPostsBO{
		UserID:  in.GetId(),
		PostIDs: in.GetPostIds(),
	}
	posts, err := s.post.ListByPostIDs(ctx, in.PostIds)
	if err != nil {
		return nil, err
	}
	if len(posts) != len(in.PostIds) {
		return nil, errorx.Err(errkey.ErrPostNotFound)
	}
	err = s.uc.AssignUserPosts(ctx, bo)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetUserPosts(ctx context.Context, in *v1.GetUserPostsRequest) (*v1.GetUserPostsReply, error) {
	postIDs, err := s.uc.GetUserPosts(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetUserPostsReply{
		PostIds: postIDs,
	}, nil
}
