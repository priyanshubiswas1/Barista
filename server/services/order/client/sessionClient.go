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
