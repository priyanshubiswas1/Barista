package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	session "github.com/priyanshubiswas1/baristacontracts/gen/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ISessionClient interface {
	GetAllSessions(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetActiveSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetAllActiveSessions(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetSessionById(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	SendSessionOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	CreateSessionAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	JoinSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type SessionClient struct {
	client  *grpc.ClientConn
	service session.SessionServiceClient
}

func NewSession(url string) (ISessionClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &SessionClient{
		client:  client,
		service: session.NewSessionServiceClient(client),
	}, nil
}

func (c *SessionClient) GetAllSessions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	shopId := params.Get("shopid")
	res, err := c.service.GetAllSessions(ctx, &session.GetAllSessionsRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	sessions := []pkg.Session{}
	for _, session := range res.Sessions {
		sessions = append(sessions, pkg.Session{
			Id:            session.SessionId,
			ShopId:        session.ShopId,
			TableId:       session.TableId,
			Otp:           session.Otp,
			Status:        session.SessionStatus,
			Amount:        session.TotalAmount,
			IsPaid:        session.IsPaid,
			PaymentMethod: session.PaymentMethod,
			PhoneNumber:   session.PhoneNumber,
			Active:        session.Active,
		})
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetAllSessionsResponse{
		Sessions: sessions,
	})
}

func (c *SessionClient) GetActiveSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	shopId := params.Get("shopid")
	tableId := params.Get("tableid")
	res, err := c.service.GetActiveSession(ctx, &session.GetActiveSessionRequest{
		ShopId:  shopId,
		TableId: tableId,
	})
	if err != nil {
		return err
	}

	// Check if session exists and is active
	if res.Session != nil && res.Session.Active {
		return pkg.WriteJSON(w, http.StatusOK, pkg.GetActiveSessionResponse{
			Session: pkg.Session{
				Id:            res.Session.SessionId,
				ShopId:        res.Session.ShopId,
				TableId:       res.Session.TableId,
				Otp:           res.Session.Otp,
				Status:        res.Session.SessionStatus,
				Amount:        res.Session.TotalAmount,
				IsPaid:        res.Session.IsPaid,
				PaymentMethod: res.Session.PaymentMethod,
				PhoneNumber:   res.Session.PhoneNumber,
				Active:        res.Session.Active,
			},
			SessionToken: res.Token, // Return token for existing session
		})
	}

	// No active session found - return empty response
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetActiveSessionResponse{
		Session:      pkg.Session{},
		SessionToken: "",
	})
}

func (c *SessionClient) GetAllActiveSessions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	shopId := params.Get("shopid")
	res, err := c.service.GetAllActiveSessions(ctx, &session.GetAllSessionsRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	sessions := []pkg.ActiveSession{}
	for _, session := range res.Sessions {
		sessions = append(sessions, pkg.ActiveSession{
			Id:            session.SessionId,
			ShopId:        session.ShopId,
			TableId:       session.TableId,
			TableNumber:   session.TableName,
			Otp:           session.Otp,
			Status:        session.SessionStatus,
			Amount:        session.TotalAmount,
			IsPaid:        session.IsPaid,
			PaymentMethod: session.PaymentMethod,
			PhoneNumber:   session.PhoneNumber,
			Active:        session.Active,
		})
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetAllActiveSessionsResponse{
		Sessions: sessions,
	})
}

func (c *SessionClient) GetSessionById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	sessionId := vars["sessionid"]
	res, err := c.service.GetSessionById(ctx, &session.GetSessionRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetSessionResponse{
		Session: pkg.Session{
			Id:            res.Session.SessionId,
			ShopId:        res.Session.ShopId,
			TableId:       res.Session.TableId,
			Otp:           res.Session.Otp,
			Status:        res.Session.SessionStatus,
			Amount:        res.Session.TotalAmount,
			IsPaid:        res.Session.IsPaid,
			PaymentMethod: res.Session.PaymentMethod,
			PhoneNumber:   res.Session.PhoneNumber,
			Active:        res.Session.Active,
		},
	})
}

func (c *SessionClient) SendSessionOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sendSessionOTP := pkg.SendSessionOTPRequest{}
	if err := json.NewDecoder(r.Body).Decode(&sendSessionOTP); err != nil {
		return err
	}
	res, err := c.service.SendSessionOTP(ctx, &session.SendSessionOTPRequest{
		ShopId:  sendSessionOTP.ShopId,
		TableId: sendSessionOTP.TableId,
		PhoneNo: sendSessionOTP.PhoneNumber,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.SendSessionOTPResponse{
		Success: res.Success,
	})
}

func (c *SessionClient) CreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	createSession := pkg.CreateSessionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createSession); err != nil {
		return err
	}
	res, err := c.service.CreateSession(ctx, &session.CreateSessionRequest{
		ShopId:        createSession.ShopId,
		TableId:       createSession.TableId,
		Otp:           createSession.Otp,
		SessionStatus: createSession.Status,
		TotalAmount:   createSession.Amount,
		IsPaid:        createSession.IsPaid,
		PaymentMethod: createSession.PaymentMethod,
		Active:        createSession.Active,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateSessionResponse{
		Session: pkg.Session{
			Id:            res.Session.SessionId,
			ShopId:        res.Session.ShopId,
			TableId:       res.Session.TableId,
			Otp:           res.Session.Otp,
			Status:        res.Session.SessionStatus,
			Amount:        res.Session.TotalAmount,
			IsPaid:        res.Session.IsPaid,
			PaymentMethod: res.Session.PaymentMethod,
			PhoneNumber:   res.Session.PhoneNumber,
			Active:        res.Session.Active,
		},
		SessionToken: res.Token,
	})
}

func (c *SessionClient) CreateSessionAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	createSession := pkg.CreateSessionAdminRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createSession); err != nil {
		return err
	}
	res, err := c.service.CreateSessionAdmin(ctx, &session.CreateSessionAdminRequest{
		ShopId:        createSession.ShopId,
		TableId:       createSession.TableId,
		SessionStatus: createSession.Status,
		TotalAmount:   createSession.Amount,
		IsPaid:        createSession.IsPaid,
		PaymentMethod: createSession.PaymentMethod,
		PhoneNumber:   createSession.PhoneNumber,
		Active:        createSession.Active,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateSessionResponse{
		Session: pkg.Session{
			Id:            res.Session.SessionId,
			ShopId:        res.Session.ShopId,
			TableId:       res.Session.TableId,
			Otp:           res.Session.Otp,
			Status:        res.Session.SessionStatus,
			Amount:        res.Session.TotalAmount,
			IsPaid:        res.Session.IsPaid,
			PaymentMethod: res.Session.PaymentMethod,
			PhoneNumber:   res.Session.PhoneNumber,
			Active:        res.Session.Active,
		},
		SessionToken: res.Token,
	})
}

func (c *SessionClient) JoinSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	joinSession := pkg.JoinSessionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&joinSession); err != nil {
		return err
	}
	res, err := c.service.JoinSession(ctx, &session.JoinSessionRequest{
		ShopId:  joinSession.ShopId,
		TableId: joinSession.TableId,
		Otp:     joinSession.Otp,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.JoinSessionResponse{
		Session: pkg.Session{
			Id:            res.Session.SessionId,
			ShopId:        res.Session.ShopId,
			TableId:       res.Session.TableId,
			Otp:           res.Session.Otp,
			Status:        res.Session.SessionStatus,
			Amount:        res.Session.TotalAmount,
			IsPaid:        res.Session.IsPaid,
			PaymentMethod: res.Session.PaymentMethod,
			PhoneNumber:   res.Session.PhoneNumber,
			Active:        res.Session.Active,
		},
		SessionToken: res.Token,
	})
}

func (c *SessionClient) UpdateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	updateSession := pkg.UpdateSessionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&updateSession); err != nil {
		return err
	}
	res, err := c.service.UpdateSession(ctx, &session.UpdateSessionRequest{
		Session: &session.Session{
			SessionId:     updateSession.Session.Id,
			ShopId:        updateSession.Session.ShopId,
			TableId:       updateSession.Session.TableId,
			Otp:           updateSession.Session.Otp,
			SessionStatus: updateSession.Session.Status,
			TotalAmount:   updateSession.Session.Amount,
			IsPaid:        updateSession.Session.IsPaid,
			PaymentMethod: updateSession.Session.PaymentMethod,
			Active:        updateSession.Session.Active,
		},
	})
	if err != nil {
		return err
	}
	if res.Success {

	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateSessionResponse{
		Success: res.Success,
	})
}
