package payment

import (
	"context"
	"fmt"
	"net"

	payment "github.com/priyanshubiswas1/baristacontracts/gen/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	payment.UnimplementedPaymentServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	payment.RegisterPaymentServiceServer(server, &gRPCserver{
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	res, err := s.service.GetPayment(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	return &payment.GetPaymentResponse{
		Payment: &payment.Payment{
			Id:            res.Id,
			SessionId:     res.SessionId,
			Amount:        res.Amount,
			PaymentMethod: res.PaymentMethod,
			PaymentStatus: res.Status,
			TransactionId: res.TransactionId,
		},
	}, nil
}

func (s *gRPCserver) CreatePayment(ctx context.Context, req *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	paymentId, err := s.service.CreatePayment(ctx, req.SessionId, req.TableId, req.ShopId, req.PhoneNumber, req.PaymentMethod, req.PaymentStatus, req.TransactionId, req.Amount)
	if err != nil {
		return nil, err
	}
	return &payment.CreatePaymentResponse{
		Payment: &payment.Payment{
			Id:            paymentId,
			SessionId:     req.SessionId,
			Amount:        req.Amount,
			PaymentMethod: req.PaymentMethod,
			PaymentStatus: req.PaymentStatus,
			TransactionId: req.TransactionId,
		},
	}, nil
}
