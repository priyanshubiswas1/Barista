package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	auth "github.com/priyanshubiswas1/barista/contracts/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IClient interface {
	GetUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type Client struct {
	client  *grpc.ClientConn
	service auth.AuthServiceClient
}

func NewClient(url string) (IClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client:  client,
		service: auth.NewAuthServiceClient(client),
	}, nil
}

func (c *Client) GetUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.GetUserReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.GetUser(ctx, &auth.GetRequest{
		Id: req.Id,
	})
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &pkg.User{
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	})
}

func (c *Client) RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.RegisterReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	user := &auth.User{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}
	res, err := c.service.RegisterUser(ctx, &auth.RegisterRequest{
		User:     user,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &pkg.RegisterRes{
		Token:   res.Token,
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	})
}

func (c *Client) LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.LoginReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.LoginUser(ctx, &auth.LoginRequest{
		Identifier: &auth.LoginRequest_Email{
			Email: req.Email,
		},
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &pkg.LoginRes{
		Token:   res.Token,
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	})
}

func (c *Client) UpdateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.UpdateUserReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	user := &auth.User{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}
	res, err := c.service.UpdateUser(ctx, &auth.UpdateRequest{
		Id:       req.Id,
		User:     user,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &pkg.UpdateUserRes{
		Token:   res.Token,
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	})
}

func (c *Client) DeleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.DeleteUserReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.DeleteUser(ctx, &auth.DeleteRequest{
		Id: req.Id,
	})
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &pkg.DeleteUserRes{
		Success: res.Success,
	})
}

func writeJSON(w http.ResponseWriter, s int, data any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(data)
}
