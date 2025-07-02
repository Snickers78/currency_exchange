package auth

import (
	"context"
	"errors"
	ssov1 "user_service/gen"
	errs "user_service/internal/domain/errors"
	"user_service/internal/lib/validate"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	Register(ctx context.Context, email, password string) (id int64, err error)
}
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if ok := validate.EmailValid(req.Email); !ok {
		return nil, status.Error(codes.InvalidArgument, "Wrong email or password")
	}

	if ok := validate.PasswordValid(req.Password); !ok {
		return nil, status.Error(codes.InvalidArgument, "Wrong email or password")
	}

	id, err := s.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errs.UserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: id,
	}, nil
}
