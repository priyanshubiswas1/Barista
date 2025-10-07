package client

import (
	"context"

	"github.com/priyanshubiswas1/barista/orderservice/pkg"
	session "github.com/priyanshubiswas1/baristacontracts/gen/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ISessionClient interface {
	UpdateSession(ctx context.Context, updateSession pkg.UpdateSessionRequest) (bool, error)
	GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error)
	CreateSessionAdmin(ctx context.Context, shopId, tableId, status, paymentMethod, phoneNumber string, totalAmount float32, isPaid bool, active bool) (pkg.Session, string, error)
}

type SessionClient struct {
	client  *grpc.ClientConn
	service session.SessionServiceClient
}

func NewSession(url string) (ISessionClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &SessionClient{
		client:  client,
		service: session.NewSessionServiceClient(client),
	}, nil
}

func (c *SessionClient) UpdateSession(ctx context.Context, updateSession pkg.UpdateSessionRequest) (bool, error) {
	res, err := c.service.UpdateSession(ctx, &session.UpdateSessionRequest{
		Session: &session.Session{
			SessionId:     updateSession.Session.Id,
			ShopId:        updateSession.Session.ShopId,
			TableId:       updateSession.Session.TableId,
			Otp:           updateSession.Session.Otp,
			SessionStatus: updateSession.Session.Status,
			TotalAmount:   updateSession.Session.Amount,
			IsPaid:        updateSession.Session.IsPaid,
			PaymentMethod: updateSession.Session.PaymentMethod,
			Active:        updateSession.Session.Active,
		},
	})
	if err != nil {
		return false, err
	}
	return res.Success, err
}

func (c *SessionClient) GetSessionById(ctx context.Context, sessionId string) (pkg.Session, error) {
	res, err := c.service.GetSessionById(ctx, &session.GetSessionRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return pkg.Session{}, err
	}
	return pkg.Session{
		Id:            res.Session.SessionId,
		ShopId:        res.Session.ShopId,
		TableId:       res.Session.TableId,
		Otp:           res.Session.Otp,
		Status:        res.Session.SessionStatus,
		Amount:        res.Session.TotalAmount,
		IsPaid:        res.Session.IsPaid,
		PaymentMethod: res.Session.PaymentMethod,
		Active:        res.Session.Active,
	}, nil
}

func (c *SessionClient) CreateSessionAdmin(ctx context.Context, shopId, tableId, status, paymentMethod, phoneNumber string, totalAmount float32, isPaid bool, active bool) (pkg.Session, string, error) {
	res, err := c.service.CreateSessionAdmin(ctx, &session.CreateSessionAdminRequest{
		ShopId:        shopId,
		TableId:       tableId,
		SessionStatus: status,
		TotalAmount:   totalAmount,
		IsPaid:        isPaid,
		PaymentMethod: paymentMethod,
		PhoneNumber:   phoneNumber,
		Active:        active,
	})
	if err != nil {
		return pkg.Session{}, "", err
	}
	return pkg.Session{
		Id:            res.Session.SessionId,
		ShopId:        res.Session.ShopId,
		TableId:       res.Session.TableId,
		Otp:           res.Session.Otp,
		Status:        res.Session.SessionStatus,
		Amount:        res.Session.TotalAmount,
		IsPaid:        res.Session.IsPaid,
		PaymentMethod: res.Session.PaymentMethod,
		PhoneNumber:   res.Session.PhoneNumber,
		Active:        res.Session.Active,
	}, res.Token, nil
}
