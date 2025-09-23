package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/sessionservice/pkg"
)

type RedisData struct {
	ShopId      string `json:"shopId"`
	TableId     string `json:"tableId"`
	PhoneNumber int32  `json:"phoneNumber"`
	OTP         string `json:"otp"`
}

type IService interface {
	GetAllSessions(ctx context.Context, shopId string) ([]pkg.Session, error)
	GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error)
	GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error)
	SendSessionOTP(ctx context.Context, shopId, tableId string, phoneNumber int32) (bool, error)
	CreateSession(ctx context.Context, shopId, tableId, status, paymentMethod, otp string, totalAmount float32, isPaid bool, active bool) (string, string, int32, error)
	JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, string, error)
	UpdateSession(ctx context.Context, id, status, paymentMethod string, totalAmount float32, isPaid bool, active bool) error
}

type service struct {
	repository  IRepository
	redisClient IRedisClient
}

func NewService(repository IRepository, redisClient IRedisClient) (IService, error) {
	return &service{
		repository:  repository,
		redisClient: redisClient,
	}, nil
}

func (s *service) GetAllSessions(ctx context.Context, shopId string) ([]pkg.Session, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	return s.repository.GetAllSessions(ctx, shopId)
}

func (s *service) GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error) {
	if shopId == "" {
		return pkg.Session{}, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return pkg.Session{}, errors.New("invalid shop ID format")
	}
	if tableId == "" {
		return pkg.Session{}, errors.New("table ID is required")
	}
	if _, err := uuid.Parse(tableId); err != nil {
		return pkg.Session{}, errors.New("invalid table ID format")
	}
	return s.repository.GetActiveSession(ctx, shopId, tableId)
}

func (s *service) GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error) {
	if sessionId == "" {
		return pkg.Session{}, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return pkg.Session{}, errors.New("invalid shop ID format")
	}
	return s.repository.GetSessionById(ctx, sessionId)
}

func (s *service) SendSessionOTP(ctx context.Context, shopId, tableId string, phoneNumber int32) (bool, error) {
	otp := generateOTP()
	value, err := json.Marshal(RedisData{
		ShopId:      shopId, //not used
		TableId:     tableId,
		PhoneNumber: phoneNumber,
		OTP:         otp,
	})
	if err != nil {
		return false, err
	}
	if err := s.redisClient.SetValue(ctx, tableId, value, time.Hour*3); err != nil {
		return false, err
	}
	log.Println("OTP: ", otp)
	return true, nil
}

func (s *service) CreateSession(ctx context.Context, shopId, tableId, status, paymentMethod, otp string, totalAmount float32, isPaid bool, active bool) (string, string, int32, error) {
	id := uuid.New().String()
	value, err := s.redisClient.GetValue(ctx, tableId)
	if err != nil {
		return "", "", 0, err
	}
	unMarshalValue := RedisData{}
	if err := json.Unmarshal([]byte(value), &unMarshalValue); err != nil {
		return "", "", 0, err
	}

	if tableId != unMarshalValue.TableId {
		return "", "", 0, errors.New("invalid table ID")
	}
	if otp != unMarshalValue.OTP {
		return "", "", 0, errors.New("invalid OTP")
	}

	phoneNumber := unMarshalValue.PhoneNumber
	if err := s.repository.CreateSession(ctx, id, shopId, tableId, status, paymentMethod, otp, phoneNumber, totalAmount, isPaid, true, time.Now()); err != nil {
		return "", "", 0, err
	}
	tokenString, err := GenerateToken(id)
	if err != nil {
		log.Println("error generating token")
		return "", "", 0, err
	}
	return id, tokenString, phoneNumber, nil
}

func (s *service) JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, string, error) {
	res, err := s.repository.JoinSession(ctx, shopId, tableId, otp)
	if err != nil {
		return pkg.Session{}, "", err
	}
	tokenString, err := GenerateToken(res.Id)
	if err != nil {
		log.Println("error generating token")
		return pkg.Session{}, "", err
	}
	return res, tokenString, nil
}

func (s *service) UpdateSession(ctx context.Context, id, status, paymentMethod string, totalAmount float32, isPaid bool, active bool) error {
	return s.repository.UpdateSession(ctx, id, status, paymentMethod, totalAmount, isPaid, active, time.Now())
}

func GenerateToken(id string) (string, error) {
	expiration := time.Now().Add(3 * time.Hour)

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

func generateOTP() string {
	otp := rand.Intn(9000) + 1000
	return fmt.Sprintf("%04d", otp) // Formats as a 4-digit string
}
