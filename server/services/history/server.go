package history

import (
	"context"
	"fmt"
	"net"

	"github.com/priyanshubiswas1/barista/historyservice/pkg"
	history "github.com/priyanshubiswas1/baristacontracts/gen/history"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCserver struct {
	service IService
	history.HistoryServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	history.RegisterHistoryServiceServer(server, &gRPCserver{
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) GetSessionSummary(ctx context.Context, req *history.SessionSummaryRequest) (*history.SessionSummaryResponse, error) {
	summ, totalAmount, totalSessions, err := s.service.GetSessionSummary(ctx, req.Offset, req.Limit, req.ShopId, req.GetTimeRange())
	if err != nil {
		return nil, err
	}
	var summaries []*history.SessionSummary
	for _, summary := range summ {
		var orderItems []*history.OrderItems
		for _, oi := range summary.OrderedItems {
			orderItems = append(orderItems, &history.OrderItems{
				OrderItemId: oi.OrderItemId,
				ItemPrice:   oi.ItemPrice,
				Quantity:    oi.Quantity,
				OrderId:     oi.OrderId,
			})
		}
		summaries = append(summaries, &history.SessionSummary{
			SessionId:     summary.SessionId,
			ShopId:        summary.ShopId,
			SessionNumber: summary.SessionNumber,
			TotalAmount:   summary.TotalAmount,
			CreatedAt:     timestamppb.New(summary.CreatedAt),
			UpdatedAt:     timestamppb.New(summary.UpdatedAt),
			OrderItems:    orderItems,
		})
	}
	return &history.SessionSummaryResponse{
		TotalAmount:    totalAmount,
		TotalSessions:  totalSessions,
		SessionSummary: summaries,
	}, nil
}

func (s *gRPCserver) CreateSession(ctx context.Context, req *history.CreateSessionRequest) (*history.SessionSummary, error) {
	id, createdAt, err := s.service.CreateSession(ctx, req.ShopId, req.SessionId, req.SessionNumber, 0, nil)
	if err != nil {
		return nil, err
	}
	return &history.SessionSummary{
		Id:            id,
		SessionId:     req.SessionId,
		ShopId:        req.ShopId,
		SessionNumber: req.SessionNumber,
		TotalAmount:   0,
		CreatedAt:     timestamppb.New(createdAt),
		UpdatedAt:     timestamppb.New(createdAt),
		OrderItems:    nil,
	}, nil
}

func (s *gRPCserver) UpdateSession(ctx context.Context, req *history.UpdateSessionRequest) (*history.SessionSummary, error) {
	id, updatedAt, err := s.service.UpdateSession(ctx, req.SessionId, req.TotalAmount)
	if err != nil {
		return nil, err
	}
	return &history.SessionSummary{
		Id:          id,
		SessionId:   req.SessionId,
		TotalAmount: req.TotalAmount,
		UpdatedAt:   timestamppb.New(updatedAt),
	}, nil
}

func (s *gRPCserver) CloseSession(ctx context.Context, req *history.CloseSessionRequest) (*history.StatusResponse, error) {
	if err := s.service.CloseSession(ctx, req.SessionId); err != nil {
		return &history.StatusResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	return &history.StatusResponse{
		Success: true,
		Message: "Session Closed successfully",
	}, nil
}

func (s *gRPCserver) AddOrder(ctx context.Context, req *history.AddOrderRequest) (*history.SessionSummary, error) {
	var orderedItems []pkg.OrderedItem
	orderItemsList := req.GetOrderItems()
	for _, oi := range orderItemsList {
		orderedItems = append(orderedItems, pkg.OrderedItem{
			OrderItemId: oi.GetOrderItemId(),
			OrderId:     oi.GetOrderId(),
			ItemPrice:   oi.GetItemPrice(),
			Quantity:    oi.GetQuantity(),
		})
	}

	summary, err := s.service.AddOrder(ctx, req.GetSessionId(), orderedItems)
	if err != nil {
		return nil, err
	}
	var orderItems []*history.OrderItems
	for _, oi := range summary.OrderedItems {
		orderItems = append(orderItems, &history.OrderItems{
			OrderItemId: oi.OrderItemId,
			ItemPrice:   oi.ItemPrice,
			Quantity:    oi.Quantity,
			OrderId:     oi.OrderId,
		})
	}
	return &history.SessionSummary{
		Id:            summary.Id,
		SessionId:     summary.SessionId,
		ShopId:        summary.ShopId,
		SessionNumber: summary.SessionNumber,
		TotalAmount:   summary.TotalAmount,
		OrderItems:    orderItems,
		CreatedAt:     timestamppb.New(summary.CreatedAt),
		UpdatedAt:     timestamppb.New(summary.UpdatedAt),
	}, nil
}

func (s *gRPCserver) RemoveOrderItems(ctx context.Context, req *history.RemoveOrderItemsRequest) (*history.StatusResponse, error) {
	err := s.service.RemoveOrderItems(ctx, req.SessionId, req.OrderItemIds)
	if err != nil {
		return &history.StatusResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	return &history.StatusResponse{
		Success: true,
		Message: "Items removed",
	}, nil
}
