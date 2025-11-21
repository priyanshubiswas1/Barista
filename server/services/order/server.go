package order

import (
	"context"
	"fmt"
	"net"

	"github.com/priyanshubiswas1/barista/orderservice/pkg"
	order "github.com/priyanshubiswas1/baristacontracts/gen/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	order.UnimplementedOrderServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	order.RegisterOrderServiceServer(server, &gRPCserver{
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) GetOrdersBySessionId(ctx context.Context, req *order.GetOrdersBySessionIdRequest) (*order.GetOrdersBySessionIdResponse, error) {
	orderItemsRes, totalAmount, err := s.service.GetOrdersBySessionId(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	orderItems := []*order.OrderItem{}
	for _, oi := range orderItemsRes {
		orderItems = append(orderItems, &order.OrderItem{
			OrderItemId:      oi.Id,
			OrderId:          oi.OrderId,
			OrderNumber:      oi.OrderNumber,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			OrderStatus:      oi.OrderStatus,
		})
	}
	return &order.GetOrdersBySessionIdResponse{
		Orderitems:        orderItems,
		TotalSessionPrice: totalAmount,
	}, nil
}

func (s *gRPCserver) GetOrderById(ctx context.Context, req *order.GetOrderByIdRequest) (*order.GetOrderByIdResponse, error) {
	orderItemsRes, totalAmount, err := s.service.GetOrderById(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	orderItems := []*order.OrderItem{}
	for _, oi := range orderItemsRes {
		orderItems = append(orderItems, &order.OrderItem{
			OrderItemId:      oi.Id,
			OrderId:          oi.OrderId,
			OrderNumber:      oi.OrderNumber,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			OrderStatus:      oi.OrderStatus,
		})
	}
	return &order.GetOrderByIdResponse{
		Orderitems:      orderItems,
		TotalOrderPrice: totalAmount,
	}, nil
}

func (s *gRPCserver) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	orderItems := []pkg.OrderItemsReq{}
	for _, oi := range req.OrderItems {
		orderItems = append(orderItems, pkg.OrderItemsReq{
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
		})
	}
	orderRes, err := s.service.CreateOrder(ctx, pkg.OrderReq{
		SessionId:    req.SessionId,
		SpecialNotes: req.SpecialNotes,
		OrderItems:   orderItems,
	})
	if err != nil {
		return nil, err
	}
	newOrderItems := []*order.OrderItem{}
	for _, oi := range orderRes.OrderItems {
		newOrderItems = append(newOrderItems, &order.OrderItem{
			OrderItemId:      oi.Id,
			OrderId:          oi.OrderId,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			OrderStatus:      oi.OrderStatus,
		})
	}
	return &order.CreateOrderResponse{
		Order: &order.Order{
			OrderId:      orderRes.Order.Id,
			SessionId:    orderRes.Order.SessionId,
			OrderStatus:  orderRes.Order.OrderStatus,
			SpecialNotes: orderRes.Order.SpecialNotes,
		},
		OrderItems: newOrderItems,
	}, nil
}

func (s *gRPCserver) CreateOrderAdmin(ctx context.Context, req *order.CreateOrderAdminRequest) (*order.CreateOrderResponse, error) {
	orderItems := []pkg.OrderItemsReq{}
	for _, oi := range req.OrderItems {
		orderItems = append(orderItems, pkg.OrderItemsReq{
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
		})
	}
	orderRes, err := s.service.CreateOrderAdmin(ctx, pkg.OrderAdminReq{
		ShopId:       req.ShopId,
		TableId:      req.TableId,
		SpecialNotes: req.SpecialNotes,
		OrderItems:   orderItems,
	})
	if err != nil {
		return nil, err
	}
	newOrderItems := []*order.OrderItem{}
	for _, oi := range orderRes.OrderItems {
		newOrderItems = append(newOrderItems, &order.OrderItem{
			OrderItemId:      oi.Id,
			OrderId:          oi.OrderId,
			ItemId:           oi.ItemId,
			Quantity:         oi.Quantity,
			PriceAtOrderTime: oi.PriceAtOrderTime,
			OrderStatus:      oi.OrderStatus,
		})
	}
	return &order.CreateOrderResponse{
		Order: &order.Order{
			OrderId:      orderRes.Order.Id,
			SessionId:    orderRes.Order.SessionId,
			OrderStatus:  orderRes.Order.OrderStatus,
			SpecialNotes: orderRes.Order.SpecialNotes,
		},
		OrderItems: newOrderItems,
	}, nil
}

func (s *gRPCserver) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	res := s.service.UpdateOrder(ctx, req.OrderId, req.OrderStatus)
	return &order.UpdateOrderResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	res := s.service.CancelOrder(ctx, req.OrderId)
	return &order.CancelOrderResponse{
		Success: res,
	}, nil
}
