package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/priyanshubiswas1/barista/apigateway/middleware"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	payment "github.com/priyanshubiswas1/baristacontracts/gen/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IPaymentClient interface {
	CreatePayment(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetPayment(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type PaymentClient struct {
	client  *grpc.ClientConn
	service payment.PaymentServiceClient
}

func NewPayment(url string) (IPaymentClient, error) {
	client, err := grpc.NewClient(pkg.StripScheme(url), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &PaymentClient{}, err
	}
	return &PaymentClient{
		client:  client,
		service: payment.NewPaymentServiceClient(client),
	}, nil
}

func (c *PaymentClient) CreatePayment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	paymentReq := pkg.PaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		return err
	}
	res, err := c.service.CreatePayment(ctx, &payment.CreatePaymentRequest{
		SessionId:     paymentReq.SessionId,
		TableId:       paymentReq.TableId,
		ShopId:        paymentReq.ShopId,
		PhoneNumber:   paymentReq.PhoneNumber,
		Amount:        paymentReq.Amount,
		PaymentMethod: paymentReq.PaymentMethod,
		PaymentStatus: paymentReq.PaymentStatus,
		TransactionId: paymentReq.TransactionId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.Payment{
		Id:            res.Payment.Id,
		SessionId:     res.Payment.SessionId,
		Amount:        res.Payment.Amount,
		PaymentMethod: res.Payment.PaymentMethod,
		PaymentStatus: res.Payment.PaymentStatus,
		TransactionId: res.Payment.TransactionId,
	})
}

func (c *PaymentClient) GetPayment(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sessIdVal := r.Context().Value(middleware.CtxSessionKeyID)
	sessIdStr, ok := sessIdVal.(string)
	if !ok || sessIdStr == "" {
		return errors.New("missing or invalid session ID in context")
	}
	res, err := c.service.GetPayment(ctx, &payment.GetPaymentRequest{
		SessionId: sessIdStr,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.Payment{
		Id:            res.Payment.Id,
		SessionId:     res.Payment.SessionId,
		Amount:        res.Payment.Amount,
		PaymentMethod: res.Payment.PaymentMethod,
		PaymentStatus: res.Payment.PaymentStatus,
		TransactionId: res.Payment.TransactionId,
	})
}
