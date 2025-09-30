package session

import (
	"context"
	"fmt"
	"net"

	session "github.com/priyanshubiswas1/baristacontracts/gen/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	session.UnimplementedSessionServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	session.RegisterSessionServiceServer(server, &gRPCserver{
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) GetAllSessions(ctx context.Context, req *session.GetAllSessionsRequest) (*session.GetAllSessionsResponse, error) {
	res, err := s.service.GetAllSessions(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	sessions := []*session.Session{}
	for _, s := range res {
		sessions = append(sessions, &session.Session{
			SessionId:     s.Id,
			ShopId:        s.ShopId,
			TableId:       s.TableId,
			Otp:           s.Otp,
			SessionStatus: s.Status,
			TotalAmount:   s.Amount,
			IsPaid:        s.IsPaid,
			PaymentMethod: s.PaymentMethod,
			PhoneNumber:   s.PhoneNumber,
			Active:        s.Active,
		})
	}
	return &session.GetAllSessionsResponse{
		Sessions: sessions,
	}, nil
}

func (s *gRPCserver) GetActiveSession(ctx context.Context, req *session.GetActiveSessionRequest) (*session.GetActiveSessionResponse, error) {
	res, err := s.service.GetActiveSession(ctx, req.ShopId, req.TableId)
	if err != nil {
		return nil, err
	}
	return &session.GetActiveSessionResponse{
		Session: &session.Session{
			SessionId:     res.Id,
			ShopId:        res.ShopId,
			TableId:       res.TableId,
			Otp:           res.Otp,
			SessionStatus: res.Status,
			TotalAmount:   res.Amount,
			IsPaid:        res.IsPaid,
			PaymentMethod: res.PaymentMethod,
			PhoneNumber:   res.PhoneNumber,
			Active:        res.Active,
		},
	}, nil
}

func (s *gRPCserver) GetSessionById(ctx context.Context, req *session.GetSessionRequest) (*session.GetSessionResponse, error) {
	res, err := s.service.GetSessionById(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	return &session.GetSessionResponse{
		Session: &session.Session{
			SessionId:     res.Id,
			ShopId:        res.ShopId,
			TableId:       res.TableId,
			Otp:           res.Otp,
			SessionStatus: res.Status,
			TotalAmount:   res.Amount,
			IsPaid:        res.IsPaid,
			PaymentMethod: res.PaymentMethod,
			PhoneNumber:   res.PhoneNumber,
			Active:        res.Active,
		},
	}, nil
}

func (s *gRPCserver) SendSessionOTP(ctx context.Context, req *session.SendSessionOTPRequest) (*session.SendSessionOTPResponse, error) {
	res, err := s.service.SendSessionOTP(ctx, req.ShopId, req.TableId, req.PhoneNo)
	if err != nil {
		return nil, err
	}
	return &session.SendSessionOTPResponse{
		Success: res,
	}, nil
}

func (s *gRPCserver) CreateSession(ctx context.Context, req *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	id, token, phoneNo, err := s.service.CreateSession(ctx, req.ShopId, req.TableId, req.SessionStatus, req.PaymentMethod, req.Otp, req.TotalAmount, req.IsPaid, req.Active)
	if err != nil {
		return nil, err
	}
	return &session.CreateSessionResponse{
		Session: &session.Session{
			SessionId:     id,
			ShopId:        req.ShopId,
			TableId:       req.TableId,
			Otp:           req.Otp,
			SessionStatus: req.SessionStatus,
			TotalAmount:   req.TotalAmount,
			IsPaid:        req.IsPaid,
			PaymentMethod: req.PaymentMethod,
			PhoneNumber:   phoneNo,
			Active:        req.Active,
		},
		Token: token,
	}, nil
}

func (s *gRPCserver) JoinSession(ctx context.Context, req *session.JoinSessionRequest) (*session.JoinSessionResponse, error) {
	resSession, tokenString, err := s.service.JoinSession(ctx, req.ShopId, req.TableId, req.Otp)
	if err != nil {
		return nil, err
	}
	return &session.JoinSessionResponse{
		Session: &session.Session{
			SessionId:     resSession.Id,
			ShopId:        resSession.ShopId,
			TableId:       resSession.TableId,
			Otp:           resSession.Otp,
			SessionStatus: resSession.Status,
			TotalAmount:   resSession.Amount,
			IsPaid:        resSession.IsPaid,
			PaymentMethod: resSession.PaymentMethod,
			PhoneNumber:   resSession.PhoneNumber,
			Active:        resSession.Active,
		},
		Token: tokenString,
	}, nil
}

func (s *gRPCserver) UpdateSession(ctx context.Context, req *session.UpdateSessionRequest) (*session.UpdateSessionResponse, error) {
	if err := s.service.UpdateSession(ctx, req.Session.SessionId, req.Session.SessionStatus, req.Session.PaymentMethod, req.Session.TotalAmount, req.Session.IsPaid, req.Session.Active); err != nil {
		return nil, err
	}
	return &session.UpdateSessionResponse{
		Success: true,
	}, nil
}
