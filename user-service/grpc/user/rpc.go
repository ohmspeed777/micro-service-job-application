package user

import (
	up "app/internal/core/proto/user"
	"app/internal/core/domain"
	"app/internal/core/services/user"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPC struct {
	user user.IUserService
	up.UnimplementedUserServer
}

func NewGrpcServer(user user.IUserService, server *grpc.Server) {
	up.RegisterUserServer(server, &GRPC{
		user: user,
	})
}

func (g *GRPC) GetProfile(ctx context.Context, req *up.GetUserReq) (*up.GetUserRes, error) {
	user, err := g.user.GetOne(ctx, domain.GetOneReq{ID: req.GetId()})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &up.GetUserRes{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
