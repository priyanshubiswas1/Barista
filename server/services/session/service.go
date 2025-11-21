package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/priyanshubiswas1/barista/sessionservice/client"
	"github.com/priyanshubiswas1/barista/sessionservice/pkg"
)

type RedisData struct {
	ShopId      string `json:"shopId"`
	TableId     string `json:"tableId"`
	PhoneNumber string `json:"phoneNumber"`
	OTP         string `json:"otp"`
}

type TableStatus string

const (
	TableStatusFree     TableStatus = "free"
	TableStatusOccupied TableStatus = "occupied"
)

type IService interface {
	GetAllSessions(ctx context.Context, shopId string) ([]pkg.Session, error)
	GetActiveSession(ctx context.Context, shopId, tableId string) (pkg.Session, error)
	GetAllActiveSessions(ctx context.Context, shopId string) ([]pkg.ActiveSession, error)
	GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error)
	SendSessionOTP(ctx context.Context, shopId, tableId, phoneNumber string) (bool, error)
	CreateSession(ctx context.Context, shopId, tableId, status, paymentMethod, otp string, totalAmount float32, isPaid bool, active bool) (string, string, string, error)
	CreateOrJoinSessionAdmin(ctx context.Context, shopId, tableId, status, paymentMethod, phoneNumber string, totalAmount float32, isPaid bool, active bool) (pkg.Session, string, error)
	JoinSession(ctx context.Context, shopId, tableId, otp string) (pkg.Session, string, error)
	UpdateSession(ctx context.Context, id, shopId, tableId, status, paymentMethod string, totalAmount float32, isPaid bool, active bool) error
}

type service struct {
	repository    IRepository
	redisClient   IRedisClient
	tableClient   client.ITableClient
	historyClient client.IHistoryClient
}

func NewService(repository IRepository, redisClient IRedisClient, tableClient client.ITableClient, historyClient client.IHistoryClient) (IService, error) {
	return &service{
		repository:    repository,
		redisClient:   redisClient,
		tableClient:   tableClient,
		historyClient: historyClient,
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
	session, err := s.repository.GetActiveSession(ctx, shopId, tableId)
	if err != nil {
		return pkg.Session{}, err
	}
	if session == (pkg.Session{}) {
		return pkg.Session{}, errors.New("no active sessions found")
	}
	return session, nil
}

func (s *service) GetAllActiveSessions(ctx context.Context, shopId string) ([]pkg.ActiveSession, error) {
	if shopId == "" {
		return nil, errors.New("shop ID is required")
	}
	if _, err := uuid.Parse(shopId); err != nil {
		return nil, errors.New("invalid shop ID format")
	}
	sessions, err := s.repository.GetAllActiveSessions(ctx, shopId)
	if err != nil {
		return nil, err
	}
	activeSessions := []pkg.ActiveSession{}
	for _, session := range sessions {
		table, err := s.tableClient.GetTable(ctx, session.TableId, session.ShopId)
		if err != nil {
			return nil, err
		}
		activeSessions = append(activeSessions, pkg.ActiveSession{
			Id:            session.Id,
			ShopId:        session.ShopId,
			TableId:       session.TableId,
			TableName:     table.Table.TableNumber,
			Otp:           session.Otp,
			Status:        session.Status,
			Amount:        session.Amount,
			IsPaid:        session.IsPaid,
			PaymentMethod: session.PaymentMethod,
			Active:        session.Active,
			PhoneNumber:   session.PhoneNumber,
		})
	}
	return activeSessions, nil
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

func (s *service) SendSessionOTP(ctx context.Context, shopId, tableId, phoneNumber string) (bool, error) {
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

func (s *service) CreateSession(ctx context.Context, shopId, tableId, status, paymentMethod, otp string, totalAmount float32, isPaid bool, active bool) (string, string, string, error) {
	id := uuid.New().String()
	value, err := s.redisClient.GetValue(ctx, tableId)
	if err != nil {
		return "", "", "", err
	}
	unMarshalValue := RedisData{}
	if err := json.Unmarshal([]byte(value), &unMarshalValue); err != nil {
		return "", "", "", err
	}

	if tableId != unMarshalValue.TableId {
		return "", "", "", errors.New("invalid table ID")
	}
	if otp != unMarshalValue.OTP {
		return "", "", "", errors.New("invalid OTP")
	}

	phoneNumber := unMarshalValue.PhoneNumber
	session, err := s.repository.GetActiveSession(ctx, shopId, tableId)
	if err != nil {
		return "", "", "", err
	}
	if session != (pkg.Session{}) {
		return "", "", "", errors.New("session already active")
	}

	sessionNoStr, err := s.repository.CreateSession(ctx, id, shopId, tableId, status, paymentMethod, otp, phoneNumber, totalAmount, isPaid, true, time.Now())
	if err != nil {
		return "", "", "", err
	}
	sessionNo, err := strconv.Atoi(sessionNoStr)
	if err != nil {
		return "", "", "", err
	}
	tokenString, err := GenerateToken(id)
	if err != nil {
		log.Println("error generating token")
		return "", "", "", err
	}
	s.redisClient.DeleteKey(ctx, tableId)

	if err := s.historyClient.CreateSession(ctx, id, shopId, int32(sessionNo)); err != nil {
		log.Printf("error creating session in history service: %v", err)
		return "", "", "", err
	}

	if res, err := s.tableClient.UpdateTableStatus(ctx, tableId, shopId, string(TableStatusOccupied)); err != nil || !res.Success {
		log.Printf("error updating table status in table service: %v", err)
		return "", "", "", err
	}
	return id, tokenString, phoneNumber, nil
}

func (s *service) CreateOrJoinSessionAdmin(ctx context.Context, shopId, tableId, status, paymentMethod, phoneNumber string, totalAmount float32, isPaid bool, active bool) (pkg.Session, string, error) {
	session, err := s.repository.GetActiveSession(ctx, shopId, tableId)
	if err != nil {
		return pkg.Session{}, "", err
	}
	if session != (pkg.Session{}) {
		res, err := s.repository.JoinSessionAdmin(ctx, shopId, tableId)
		if err != nil {
			return pkg.Session{}, "", err
		}
		tokenString, err := GenerateToken(res.Id)
		if err != nil {
			log.Println("error generating token")
			return pkg.Session{}, tokenString, err
		}
		return res, tokenString, nil
	} else {
		otp := generateOTP()
		//create session
		id := uuid.New().String()
		sessionNoStr, err := s.repository.CreateSession(ctx, id, shopId, tableId, status, paymentMethod, otp, phoneNumber, totalAmount, isPaid, true, time.Now())
		if err != nil {
			return pkg.Session{}, "", err
		}
		sessionNo, err := strconv.Atoi(sessionNoStr)
		if err != nil {
			return pkg.Session{}, "", err
		}
		tokenString, err := GenerateToken(id)
		if err != nil {
			log.Println("error generating token")
			return pkg.Session{}, tokenString, err
		}

		if err := s.historyClient.CreateSession(ctx, id, shopId, int32(sessionNo)); err != nil {
			log.Printf("error creating session in history service: %v", err)
			return pkg.Session{}, "", err
		}

		if res, err := s.tableClient.UpdateTableStatus(ctx, tableId, shopId, string(TableStatusOccupied)); err != nil || !res.Success {
			log.Printf("error updating table status in table service: %v", err)
			return pkg.Session{}, "", err
		}

		return pkg.Session{
			Id:            id,
			ShopId:        shopId,
			TableId:       tableId,
			Otp:           otp,
			Status:        status,
			Amount:        totalAmount,
			IsPaid:        isPaid,
			PaymentMethod: paymentMethod,
			Active:        active,
			PhoneNumber:   phoneNumber,
		}, tokenString, nil
	}
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

func (s *service) UpdateSession(ctx context.Context, id, shopId, tableId, status, paymentMethod string, totalAmount float32, isPaid bool, active bool) error {
	err := s.repository.UpdateSession(ctx, id, status, paymentMethod, totalAmount, isPaid, active, time.Now())
	if err != nil {
		return err
	}
	if !active {
		err := s.historyClient.CloseSession(ctx, id)
		if err != nil {
			return err
		}
		if res, err := s.tableClient.UpdateTableStatus(ctx, tableId, shopId, string(TableStatusFree)); err != nil || !res.Success {
			log.Printf("error updating table status in table service: %v", err)
			return err
		}
	}
	// else {
	// 	err := s.historyClient.UpdateSession(ctx, id, totalAmount)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
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
