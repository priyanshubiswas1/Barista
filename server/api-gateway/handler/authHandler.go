package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	auth "github.com/priyanshubiswas1/baristacontracts/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IAuthClient interface {
	GetUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type AuthClient struct {
	client  *grpc.ClientConn
	service auth.AuthServiceClient
}

func NewClient(url string) (IAuthClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		client:  client,
		service: auth.NewAuthServiceClient(client),
	}, nil
}

func (c *AuthClient) GetUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
	res, err := c.service.GetUser(ctx, &auth.GetRequest{
		Id: idStr,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.User{
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	})
}

func (c *AuthClient) RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
	return pkg.WriteJSON(w, http.StatusOK, &pkg.RegisterRes{
		Token:   res.Token,
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
		ShopId:  res.ShopId,
		MenuId:  res.MenuId,
	})
}

func (c *AuthClient) LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
	return pkg.WriteJSON(w, http.StatusOK, &pkg.LoginRes{
		Token:   res.Token,
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
		ShopId:  res.ShopId,
		MenuId:  res.MenuId,
	})
}

func (c *AuthClient) UpdateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
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
		Id:          idStr,
		User:        user,
		Password:    req.Password,
		Oldpassword: req.OldPassword,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.UpdateUserRes{
		Token:   res.Token,
		Name:    res.User.Name,
		Email:   res.User.Email,
		Phone:   res.User.Phone,
		Address: res.User.Address,
	})
}

func (c *AuthClient) DeleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}
	req := pkg.DeleteUserReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.DeleteUser(ctx, &auth.DeleteRequest{
		Id:       idStr,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, &pkg.DeleteUserRes{
		Success: res.Success,
	})
}
