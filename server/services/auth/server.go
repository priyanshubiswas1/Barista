package auth

import (
	"context"
	"fmt"
	"net"

	"github.com/priyanshubiswas1/authservice/pkg"
	auth "github.com/priyanshubiswas1/baristacontracts/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	auth.UnimplementedAuthServiceServer
}

func ListenGRPC(s Service, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	auth.RegisterAuthServiceServer(server, &grpcServer{service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) RegisterUser(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	user, err := s.service.Register(ctx, pkg.User{
		Name:     req.User.Name,
		Email:    req.User.Email,
		Password: req.Password,
		Phone:    req.User.Phone,
		Address:  req.User.Address,
	})
	if err != nil {
		return nil, err
	}

	authUser := &auth.User{
		Name:    req.User.Name,
		Email:   req.User.Email,
		Phone:   req.User.Phone,
		Address: req.User.Address,
	}

	return &auth.RegisterResponse{
		User:  authUser,
		Token: user.Token,
	}, nil
}
func (s *grpcServer) LoginUser(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return nil, nil
}
func (s *grpcServer) UpdateUser(ctx context.Context, req *auth.UpdateRequest) (*auth.UpdateResponse, error) {
	return nil, nil
}
func (s *grpcServer) DeleteUser(ctx context.Context, req *auth.DeleteRequest) (*auth.DeleteResponse, error) {
	return nil, nil
}
func (s *grpcServer) GetUser(ctx context.Context, req *auth.GetRequest) (*auth.GetResponse, error) {
	return nil, nil
}
