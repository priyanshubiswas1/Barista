package payment

import (
	"context"
	"fmt"
	"net"

	"net/http"

	paymentPb "github.com/priyanshubiswas1/baristacontracts/gen/payment"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	paymentPb.UnimplementedPaymentServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}

	m := cmux.New(lis)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := grpc.NewServer()
	paymentPb.RegisterPaymentServiceServer(grpcS, &gRPCserver{
		service: service,
	})
	reflection.Register(grpcS)

	httpS := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
			if r.URL.Path == "/ready" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Ready"))
				return
			}
			http.NotFound(w, r)
		}),
	}

	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	fmt.Printf("payment service listening on %s (gRPC and HTTP health checks)\n", port)
	return m.Serve()
}

func (s *gRPCserver) GetPayment(ctx context.Context, req *paymentPb.GetPaymentRequest) (*paymentPb.GetPaymentResponse, error) {
	res, err := s.service.GetPayment(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	return &paymentPb.GetPaymentResponse{
		Payment: &paymentPb.Payment{
			Id:            res.Id,
			SessionId:     res.SessionId,
			Amount:        res.Amount,
			PaymentMethod: res.PaymentMethod,
			PaymentStatus: res.Status,
			TransactionId: res.TransactionId,
		},
	}, nil
}

func (s *gRPCserver) CreatePayment(ctx context.Context, req *paymentPb.CreatePaymentRequest) (*paymentPb.CreatePaymentResponse, error) {
	paymentId, err := s.service.CreatePayment(ctx, req.SessionId, req.TableId, req.ShopId, req.PhoneNumber, req.PaymentMethod, req.PaymentStatus, req.TransactionId, req.Amount)
	if err != nil {
		return nil, err
	}
	return &paymentPb.CreatePaymentResponse{
		Payment: &paymentPb.Payment{
			Id:            paymentId,
			SessionId:     req.SessionId,
			Amount:        req.Amount,
			PaymentMethod: req.PaymentMethod,
			PaymentStatus: req.PaymentStatus,
			TransactionId: req.TransactionId,
		},
	}, nil
}
