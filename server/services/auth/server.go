package auth

import (
	"context"
	"errors"
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
	res, err := s.service.Register(ctx, pkg.User{
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
		Token: res.Token,
	}, nil
}
func (s *grpcServer) LoginUser(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var email, phone string
	switch id := req.Identifier.(type) {
	case *auth.LoginRequest_Email:
		email = id.Email
	case *auth.LoginRequest_Phonenumber:
		phone = id.Phonenumber
	default:
		return nil, errors.New("either email or phone number must be provided")
	}
	res, err := s.service.Login(ctx, email, phone, req.Password)
	if err != nil {
		return nil, err
	}
	authUser := &auth.User{
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	}
	return &auth.LoginResponse{
		Token: res.Token,
		User:  authUser,
	}, nil
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
