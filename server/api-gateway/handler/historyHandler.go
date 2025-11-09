package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	history "github.com/priyanshubiswas1/baristacontracts/gen/history"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IHistoryClient interface {
	GetSessionSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateSessionSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CloseSessionSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	AddOrderSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	RemoveOrderSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type HistoryClient struct {
	client  *grpc.ClientConn
	service history.HistoryServiceClient
}

func NewHistory(url string) (IHistoryClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &HistoryClient{
		client:  client,
		service: history.NewHistoryServiceClient(client),
	}, nil
}

func (c *HistoryClient) GetSessionSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	shopId := params.Get("shopid")
	offset := params.Get("offset")
	limit := params.Get("limit")

	offsetNum, err := strconv.Atoi(offset)
	if err != nil {
		return err
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	res, err := c.service.GetSessionSummary(ctx, &history.SessionSummaryRequest{
		ShopId: shopId,
		Offset: int32(offsetNum),
		Limit:  int32(limitNum),
	})
	if err != nil {
		return err
	}

	summaries := []pkg.SessionSummary{}
	for _, summary := range res.SessionSummary {
		orderedItems := []pkg.OrderedItems{}
		for _, oi := range summary.OrderItems {
			orderedItems = append(orderedItems, pkg.OrderedItems{
				OrderItemId: oi.OrderItemId,
				ItemPrice:   oi.ItemPrice,
				Quantity:    oi.Quantity,
			})
		}
		summaries = append(summaries, pkg.SessionSummary{
			SessionId:     summary.SessionId,
			ShopId:        summary.ShopId,
			SessionNumber: summary.SessionNumber,
			TotalAmount:   summary.TotalAmount,
			CreatedAt:     summary.CreatedAt.AsTime(),
			UpdatedAt:     summary.UpdatedAt.AsTime(),
			OrderedItems:  orderedItems,
		})
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.GetSessionSummaryResponse{
		SessionSummaries: summaries,
	})
}

func (c *HistoryClient) CreateSessionSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := pkg.CreateSessionSummaryRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}

	res, err := c.service.CreateSession(ctx, &history.CreateSessionRequest{
		SessionId:     payload.SessionId,
		ShopId:        payload.ShopId,
		SessionNumber: payload.SessionNumber,
	})
	if err != nil {
		return err
	}

	orderedItems := []pkg.OrderedItems{}
	for _, oi := range res.OrderItems {
		orderedItems = append(orderedItems, pkg.OrderedItems{
			OrderItemId: oi.OrderItemId,
			ItemPrice:   oi.ItemPrice,
			Quantity:    oi.Quantity,
		})
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateSessionSummaryResponse{
		SessionSummary: pkg.SessionSummary{
			SessionId:     res.SessionId,
			ShopId:        res.ShopId,
			SessionNumber: res.SessionNumber,
			TotalAmount:   res.TotalAmount,
			CreatedAt:     res.CreatedAt.AsTime(),
			UpdatedAt:     res.UpdatedAt.AsTime(),
			OrderedItems:  orderedItems,
		},
	})
}

func (c *HistoryClient) CloseSessionSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := pkg.CloseSessionSummaryRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}

	res, err := c.service.CloseSession(ctx, &history.CloseSessionRequest{
		SessionId: payload.SessionId,
	})
	if err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.StatusRespone{
		Success: res.Success,
		Message: res.Message,
	})
}

func (c *HistoryClient) AddOrderSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := pkg.AddOrderSummaryRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}

	orderItems := []*history.OrderItems{}
	for _, oi := range payload.OrderedItems {
		orderItems = append(orderItems, &history.OrderItems{
			OrderItemId: oi.OrderItemId,
			ItemPrice:   oi.ItemPrice,
			Quantity:    oi.Quantity,
		})
	}

	res, err := c.service.AddOrder(ctx, &history.AddOrderRequest{
		SessionId:  payload.SessionId,
		OrderItems: orderItems,
	})
	if err != nil {
		return err
	}

	orderedItems := []pkg.OrderedItems{}
	for _, oi := range res.OrderItems {
		orderedItems = append(orderedItems, pkg.OrderedItems{
			OrderItemId: oi.OrderItemId,
			ItemPrice:   oi.ItemPrice,
			Quantity:    oi.Quantity,
		})
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.SessionSummary{
		SessionId:     res.SessionId,
		ShopId:        res.ShopId,
		SessionNumber: res.SessionNumber,
		TotalAmount:   res.TotalAmount,
		CreatedAt:     res.CreatedAt.AsTime(),
		UpdatedAt:     res.UpdatedAt.AsTime(),
		OrderedItems:  orderedItems,
	})
}

func (c *HistoryClient) RemoveOrderSummary(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := pkg.RemoveOrderSummaryRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}

	orderItemIds := []string{}
	orderItemIds = append(orderItemIds, payload.OrderItemIds...)

	res, err := c.service.RemoveOrderItems(ctx, &history.RemoveOrderItemsRequest{
		SessionId:   payload.SessionId,
		OrderItemId: orderItemIds,
	})
	if err != nil {
		return err
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.StatusRespone{
		Success: res.Success,
		Message: res.Message,
	})
}
