package client

import (
	"context"
	"log"

	"errors"

	"github.com/priyanshubiswas1/barista/orderservice/pkg"
	history "github.com/priyanshubiswas1/baristacontracts/gen/history"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IHistoryClient interface {
	AddOrder(ctx context.Context, sessionId string, orderItems []pkg.AddOrderHistoryReq) error
	RemoveOrderItems(ctx context.Context, sessionId string, orderItemIds []string) error
}

type historyClient struct {
	client  *grpc.ClientConn
	service history.HistoryServiceClient
}

func NewHistory(url string) (IHistoryClient, error) {
	if url == "" {
		return &historyClient{}, nil
	}
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &historyClient{
		client:  client,
		service: history.NewHistoryServiceClient(client),
	}, nil
}

func (c *historyClient) AddOrder(ctx context.Context, sessionId string, orderItems []pkg.AddOrderHistoryReq) error {
	if c.service == nil {
		return errors.New("history service URL not configured")
	}
	var hisOrdItems []*history.OrderItems
	for _, oi := range orderItems {
		hisOrdItems = append(hisOrdItems, &history.OrderItems{
			OrderItemId: oi.OrderItemId,
			OrderId:     oi.OrderId,
			ItemPrice:   oi.ItemPrice,
			Quantity:    oi.Quantity,
		})
	}
	log.Printf("Client: Sending AddOrder request - SessionId: %s, OrderItems count: %d", sessionId, len(hisOrdItems))
	_, err := c.service.AddOrder(ctx, &history.AddOrderRequest{
		SessionId:  sessionId,
		OrderItems: hisOrdItems,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *historyClient) RemoveOrderItems(ctx context.Context, sessionId string, orderItemIds []string) error {
	if c.service == nil {
		return errors.New("history service URL not configured")
	}
	_, err := c.service.RemoveOrderItems(ctx, &history.RemoveOrderItemsRequest{
		SessionId:    sessionId,
		OrderItemIds: orderItemIds,
	})
	if err != nil {
		return err
	}
	return nil
}
