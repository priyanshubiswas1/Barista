package client

import (
	"context"

	history "github.com/priyanshubiswas1/baristacontracts/gen/history"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IHistoryClient interface {
	CreateSession(ctx context.Context, sessionId, shopId string, sessionNumber int32) error
	CloseSession(ctx context.Context, sessionId string) error
	UpdateSession(ctx context.Context, sessionId string, totalAmount float32) error
}

type historyClient struct {
	client  *grpc.ClientConn
	service history.HistoryServiceClient
}

func NewHistory(url string) (IHistoryClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &historyClient{
		client:  client,
		service: history.NewHistoryServiceClient(client),
	}, nil
}

func (c *historyClient) CreateSession(ctx context.Context, sessionId, shopId string, sessionNumber int32) error {
	_, err := c.service.CreateSession(ctx, &history.CreateSessionRequest{
		SessionId:     sessionId,
		ShopId:        shopId,
		SessionNumber: sessionNumber,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *historyClient) CloseSession(ctx context.Context, sessionId string) error {
	_, err := c.service.CloseSession(ctx, &history.CloseSessionRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *historyClient) UpdateSession(ctx context.Context, sessionId string, totalAmount float32) error {
	_, err := c.service.UpdateSession(ctx, &history.UpdateSessionRequest{
		SessionId:   sessionId,
		TotalAmount: totalAmount,
	})
	if err != nil {
		return err
	}
	return nil
}

// func (c *TableClient) GetTable(ctx context.Context, tableId, shopId string) (pkg.GetTableResponse, error) {
// 	res, err := c.service.GetTable(ctx, &table.GetTableRequest{
// 		Id:     tableId,
// 		ShopId: shopId,
// 	})
// 	if err != nil {
// 		return pkg.GetTableResponse{}, err
// 	}
// 	return pkg.GetTableResponse{
// 		Table: pkg.Table{
// 			Id:          res.Table.Id,
// 			ShopId:      res.Table.ShopId,
// 			TableNumber: res.Table.TableNumber,
// 			QrUrl:       res.Table.Qr,
// 		},
// 	}, nil
// }
