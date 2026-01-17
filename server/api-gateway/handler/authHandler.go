package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
)

type IAuthClient interface {
	GetUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	RegisterWithShopMenuHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type AuthClient struct {
	url    string
	client *http.Client
}

func NewClient(url string) (IAuthClient, error) {
	if url == "" {
		return &AuthClient{}, nil
	}
	// Trim trailing slash if present?
	// Usually better to be safe, but let's assume valid URL
	return &AuthClient{
		url: url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *AuthClient) GetUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/user/%s", c.url, idStr), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var res pkg.User
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, &res)
}

func (c *AuthClient) RegisterWithShopMenuHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Construct payload for auth service
	payload := map[string]any{
		"user": map[string]string{
			"name":    req.Name,
			"email":   req.Email,
			"phone":   req.Phone,
			"address": req.Address,
		},
		"password": req.Password,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/register-owner", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If error status, forward it? Or generic error?
	if resp.StatusCode != http.StatusOK {
		// Read body message
		// For now generic error
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var res pkg.RegisterUserWithShopMenuRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	// Add user details to response if missing from auth service response?
	// Auth service returns Token, ShopId, MenuId.
	// The api-gateway response struct also has Name, Email etc.
	// But `auth` service `RegisterUserWithShopMenuResponse` only had Token, ShopId, MenuId!
	// Wait, check `auth/server.go`:
	// `RegisterUserWithShopMenuResponse` struct has `User` struct inside?
	// `auth` service `server.go` populated `User`.
	// But `auth` service `pkg/models.go` `RegisterWithShopMenuResponse` struct defined as:
	// type RegisterWithShopMenuResponse struct { Token, ShopId, MenuId }
	// It did NOT have User!
	// Oh... `grpcServer` returned `auth.RegisterUserWithShopMenuResponse` (proto) which had User.
	// But `s.service.RegisterWithShopMenu` return `pkg.RegisterWithShopMenuResponse` which had NO user.
	// In `auth/server.go`, the handler manually constructed the Proto User from input req.
	// In `httpServer` (my new code), `RegisterUserWithShopMenu` returns `pkg.RegisterWithShopMenuResponse`.
	// Which means passing back User details is lost if I don't add them.
	// User details are useful for frontend maybe.
	// I should probably return User details in `pkg.RegisterWithShopMenuResponse` or just fill them from request here.
	// Given I have them in `req` here, I can fill them.

	res.Name = req.Name
	res.Email = req.Email
	res.Phone = req.Phone
	res.Address = req.Address

	return pkg.WriteJSON(w, http.StatusOK, &res)
}

func (c *AuthClient) RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	payload := map[string]any{
		"user": map[string]string{
			"name":    req.Name,
			"email":   req.Email,
			"phone":   req.Phone,
			"address": req.Address,
		},
		"password": req.Password,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/register", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	// Auth service returns Token.
	var res map[string]string // { "token": "..." }
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	response := &pkg.RegisterRes{
		Token:   res["token"],
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	return pkg.WriteJSON(w, http.StatusOK, response)
}

func (c *AuthClient) LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req pkg.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	payload := map[string]string{
		"email":       req.Email,
		"phonenumber": req.Phone,
		"password":    req.Password,
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/login", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	// auth service LoginResponse structure: Token, User, ShopId, MenuId.
	// pkg.LoginRes structure: Token, Name, Email... (flat)
	// I need to decode into intermediate struct or map, or update pkg.LoginRes
	// api-gateway `LoginRes` is flat.
	// auth service returns:
	/*
		type LoginResponse struct {
			Token  string `json:"token"`
			User   User   `json:"user"`
			ShopId string `json:"shopId"`
			MenuId string `json:"menuId"`
		}
	*/

	var authRes struct {
		Token  string   `json:"token"`
		User   pkg.User `json:"user"`
		ShopId string   `json:"shopId"`
		MenuId string   `json:"menuId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authRes); err != nil {
		return err
	}

	finalRes := &pkg.LoginRes{
		Token:   authRes.Token,
		Name:    authRes.User.Name,
		Email:   authRes.User.Email,
		Phone:   authRes.User.Phone,
		Address: authRes.User.Address,
		ShopId:  authRes.ShopId,
		MenuId:  authRes.MenuId,
	}

	return pkg.WriteJSON(w, http.StatusOK, finalRes)
}

func (c *AuthClient) UpdateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	var req pkg.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Auth service UpdateReq: {id, oldpassword, newpassword, user: {name, email...}}
	payload := map[string]any{
		"id":          idStr,
		"oldpassword": req.OldPassword,
		"newpassword": req.Password,
		"user": map[string]string{
			"name":    req.Name,
			"email":   req.Email,
			"phone":   req.Phone,
			"address": req.Address,
		},
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/user", c.url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var res map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	// Construct response.
	// Assuming successful update returns token and user?
	// I'll return what I have.
	token, _ := res["token"].(string)

	finalRes := &pkg.UpdateUserRes{
		Token:   token,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	return pkg.WriteJSON(w, http.StatusOK, finalRes)
}

func (c *AuthClient) DeleteUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idVal := r.Context().Value(middleware.CtxKeyID)
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		return errors.New("missing or invalid user ID in context")
	}

	var req pkg.DeleteUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	payload := map[string]string{
		"password": req.Password,
	}
	body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/user/%s", c.url, idStr), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	return pkg.WriteJSON(w, http.StatusOK, &pkg.DeleteUserRes{
		Success: true,
	})
}
