package auth

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/authservice/client"
	"github.com/priyanshubiswas1/barista/authservice/pkg"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, user *pkg.User) (*pkg.RegisterResponse, error)
	Login(ctx context.Context, email, phone, password string) (*pkg.LoginResponse, error)
	GetUser(ctx context.Context, id string) (*pkg.User, error)
	UpdateUser(ctx context.Context, id, password string, user *pkg.User) (*pkg.LoginResponse, error)
	DeleteUser(ctx context.Context, id, password string) error
}

type authService struct {
	repository Repository
	shopClient client.IShopClient
	menuClient client.IMenuClient
}

func NewService(r Repository, shopClient client.IShopClient, menuClient client.IMenuClient) Service {
	return &authService{
		repository: r,
		shopClient: shopClient,
		menuClient: menuClient,
	}
}

func (s *authService) Register(ctx context.Context, user *pkg.User) (*pkg.RegisterResponse, error) {
	createdAt := time.Now()
	id := uuid.New().String()
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := s.repository.Register(ctx, id, user, true, createdAt); err != nil {
		return nil, err
	}

	tokenString, err := GenerateToken(id)
	if err != nil {
		log.Println("error generating token")
		return nil, err
	}

	shopId, err := s.shopClient.RegisterShop(ctx, pkg.CreateShopReq{
		OwnerId:  id,
		ShopName: user.Name,
	})
	if err != nil {
		log.Println("error creating shop")
		return nil, err
	}

	menuId, err := s.menuClient.CreateMenu(ctx, pkg.CreateMenuReq{
		ShopId:   shopId,
		MenuName: user.Name,
	})
	if err != nil {
		log.Println("error creating menu")
		return nil, err
	}

	return &pkg.RegisterResponse{
		Token:  tokenString,
		ShopId: shopId,
		MenuId: menuId,
	}, nil
}

func (s *authService) Login(ctx context.Context, email, phone, password string) (*pkg.LoginResponse, error) {
	id, user, err := s.repository.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	tokenString, err := GenerateToken(id)
	if err != nil {
		log.Println("error generating token")
		return nil, err
	}

	shopIds, err := s.shopClient.GetShopByOwnerId(ctx, id)
	if err != nil {
		log.Println("error fetching shop")
		return nil, err
	}
	if len(shopIds) < 1 {
		log.Println("user has no shops")
	}

	menuIds, err := s.menuClient.GetMenus(ctx, shopIds[0])
	if err != nil {
		log.Println("error fetching menus")
		return nil, err
	}
	if len(menuIds) < 1 {
		log.Println("shop has no menu")
	}

	return &pkg.LoginResponse{
		Token:  tokenString,
		User:   *user,
		ShopId: shopIds[0],
		MenuId: menuIds[0],
	}, nil
}

func (s *authService) GetUser(ctx context.Context, id string) (*pkg.User, error) {
	return s.repository.Getuser(ctx, id)
}

func (s *authService) UpdateUser(ctx context.Context, id, oldPassword string, user *pkg.User) (*pkg.LoginResponse, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	updateUser, err := s.repository.UpdateUser(ctx, id, oldPassword, user)
	if err != nil {
		return nil, err
	}
	tokenString, err := GenerateToken(id)
	if err != nil {
		log.Println("error generating token")
		return nil, err
	}
	return &pkg.LoginResponse{
		Token: tokenString,
		User:  *updateUser,
	}, nil
}

func (s *authService) DeleteUser(ctx context.Context, id, password string) error {
	return s.repository.DeleteUser(ctx, id, password)
}

func GenerateToken(id string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
