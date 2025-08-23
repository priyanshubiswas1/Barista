package auth

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/priyanshubiswas1/authservice/pkg"
)

type Service interface {
	Register(ctx context.Context, user pkg.User) (*pkg.RegisterResponse, error)
	Login(ctx context.Context, email, phone, password string) (*pkg.LoginResponse, error)
	GetUser(ctx context.Context, id string) (*pkg.User, error)
	UpdateUser(ctx context.Context) error
	DeleteUser(ctx context.Context) error
}

type Auth struct {
}

type authService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &authService{r}
}

func (s *authService) Register(ctx context.Context, user pkg.User) (*pkg.RegisterResponse, error) {
	createdAt := time.Now()
	id := uuid.New().String()
	if err := s.repository.Register(ctx, id, user, true, createdAt); err != nil {
		return nil, err
	}

	tokenString, err := GenerateToken(id, user.Password)
	if err != nil {
		log.Println("error generating token")
		return nil, err
	}

	return &pkg.RegisterResponse{
		Token: tokenString,
	}, nil
}

func (s *authService) Login(ctx context.Context, email, phone, password string) (*pkg.LoginResponse, error) {
	id, user, err := s.repository.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	tokenString, err := GenerateToken(id, user.Password)
	if err != nil {
		log.Println("error generating token")
		return nil, err
	}
	return &pkg.LoginResponse{
		Token: tokenString,
		User:  *user,
	}, nil
}

func (s *authService) GetUser(ctx context.Context, id string) (*pkg.User, error) {
	return s.repository.Getuser(ctx, id)
}

func (s *authService) UpdateUser(ctx context.Context) error {
	return nil
}

func (s *authService) DeleteUser(ctx context.Context) error {
	return nil
}

func GenerateToken(id, password string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   password,
		ID:        id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("key"))
	if err != nil {
		log.Println("error while signing : ", err)
		return "", err
	}
	return signedToken, nil
}
