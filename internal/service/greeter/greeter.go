package greeter

import (
	"context"
	v1 "quest-admin/api/gen/helloworld/v1"
	"quest-admin/internal/biz/greeter"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer
	uc *greeter.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *greeter.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g := &greeter.Greeter{Hello: in.Name}
	_, err := s.uc.CreateGreeter(ctx, g)
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + in.Name}, nil
}
