package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	order "github.com/priyanshubiswas1/baristacontracts/gen/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IOrderClient interface {
	GetOrdersBySessionId(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetOrderById(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CancelOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type OrderClient struct {
	client  *grpc.ClientConn
	service order.OrderServiceClient
}

func NewOrder(url string) (IOrderClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &OrderClient{}, err
	}
	return &OrderClient{
		client:  client,
		service: order.NewOrderServiceClient(client),
	}, nil
}

func (c *OrderClient) GetOrdersBySessionId(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var sessionId string

	if idVal := r.Context().Value(middleware.CtxKeyID); idVal != nil {
		if idStr, ok := idVal.(string); ok && idStr != "" {
			vars := mux.Vars(r)
			sessionId = vars["sessionid"]
			if _, err := uuid.Parse(sessionId); err != nil {
				return errors.New("invalid session ID format")
			}
		}
	}

	if sessionId == "" {
		sessIdVal := r.Context().Value(middleware.CtxSessionKeyID)
		sessIdStr, ok := sessIdVal.(string)
		if !ok || sessIdStr == "" {
			return errors.New("missing or invalid session ID in context")
		}
		sessionId = sessIdStr
	}

	res, err := c.service.GetOrdersBySessionId(ctx, &order.GetOrdersBySessionIdRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return err
	}

	orderItems := []pkg.OrderItems{}
	for _, oi := range res.Orderitems {
		orderItems = append(orderItems, pkg.OrderItems{
			Id:               oi.OrderItemId,
			OrderId:          oi.OrderId,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			Status:           oi.OrderStatus,
		})
	}

	return pkg.WriteJSON(w, http.StatusOK, pkg.GetOrdersBySessionIdResponse{
		OrderItems: orderItems,
		TotalPrice: res.TotalSessionPrice,
	})
}

func (c *OrderClient) GetOrderById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	orderId := vars["orderid"]
	if _, err := uuid.Parse(orderId); err != nil {
		return errors.New("invalid table ID format")
	}
	res, err := c.service.GetOrderById(ctx, &order.GetOrderByIdRequest{
		OrderId: orderId,
	})
	if err != nil {
		return err
	}
	orderItems := []pkg.OrderItems{}
	for _, oi := range res.Orderitems {
		orderItems = append(orderItems, pkg.OrderItems{
			Id:               oi.OrderItemId,
			OrderId:          oi.OrderId,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			Status:           oi.OrderStatus,
		})
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetOrdersByOrderIdResponse{
		OrderItems: orderItems,
		TotalPrice: res.TotalOrderPrice,
	})
}

func (c *OrderClient) CreateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sessIdVal := r.Context().Value(middleware.CtxSessionKeyID)
	sessIdStr, ok := sessIdVal.(string)
	if !ok || sessIdStr == "" {
		return errors.New("missing or invalid session ID in context")
	}
	req := pkg.CreateOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	orderItems := []*order.OrderItemReq{}
	for _, oi := range req.OrderItems {
		orderItems = append(orderItems, &order.OrderItemReq{
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
		})
	}
	res, err := c.service.CreateOrder(ctx, &order.CreateOrderRequest{
		SessionId:    sessIdStr,
		SpecialNotes: req.SpecialNotes,
		OrderItems:   orderItems,
	})
	if err != nil {
		return err
	}
	orderItemsRes := []pkg.OrderItems{}
	for _, oi := range res.OrderItems {
		orderItemsRes = append(orderItemsRes, pkg.OrderItems{
			Id:               oi.OrderItemId,
			OrderId:          oi.OrderId,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
		})
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateOrderResponse{
		Order: pkg.Order{
			Id:           res.Order.OrderId,
			SessionId:    res.Order.SessionId,
			OrderStatus:  res.Order.OrderStatus,
			SpecialNotes: res.Order.SpecialNotes,
		},
		OrderItems: orderItemsRes,
	})
}

func (c *OrderClient) UpdateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.UpdateOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.UpdateOrder(ctx, &order.UpdateOrderRequest{
		OrderId:     req.OrderId,
		OrderStatus: req.OrderStatus,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateOrderResponse{
		Success: res.Success,
	})
}

func (c *OrderClient) CancelOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.CancelOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CancelOrder(ctx, &order.CancelOrderRequest{
		OrderId: req.OrderId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteTableResponse{
		Success: res.Success,
	})
}
