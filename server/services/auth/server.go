package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/authservice/pkg"
)

type httpServer struct {
	service Service
}

func ListenHTTP(s Service, port string) error {
	router := mux.NewRouter()
	server := &httpServer{service: s}

	// Health checks
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	}).Methods("GET")

	router.HandleFunc("/register", server.RegisterUser).Methods("POST")
	router.HandleFunc("/register-owner", server.RegisterUserWithShopMenu).Methods("POST")
	router.HandleFunc("/login", server.LoginUser).Methods("POST")
	router.HandleFunc("/user", server.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", server.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", server.DeleteUser).Methods("DELETE")

	log.Printf("auth service listening on HTTP %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func (s *httpServer) RegisterUserWithShopMenu(w http.ResponseWriter, r *http.Request) {
	var req pkg.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Map to pkg.User. Password is in req.Password.
	user := &pkg.User{
		Name:     req.User.Name,
		Email:    req.User.Email,
		Phone:    req.User.Phone,
		Address:  req.User.Address,
		Password: req.Password,
	}

	res, err := s.service.RegisterWithShopMenu(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response structure from existing code
	response := pkg.RegisterWithShopMenuResponse{
		Token:  res.Token,
		ShopId: res.ShopId,
		MenuId: res.MenuId,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *httpServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req pkg.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &pkg.User{
		Name:     req.User.Name,
		Email:    req.User.Email,
		Phone:    req.User.Phone,
		Address:  req.User.Address,
		Password: req.Password,
	}

	res, err := s.service.Register(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := pkg.RegisterResponse{
		Token: res.Token,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *httpServer) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req pkg.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := s.service.Login(r.Context(), req.Email, req.Phonenumber, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := pkg.LoginResponse{
		Token:  res.Token,
		User:   res.User,
		ShopId: res.ShopId,
		MenuId: res.MenuId,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *httpServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req pkg.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &pkg.User{
		Name:     req.User.Name,
		Email:    req.User.Email,
		Phone:    req.User.Phone,
		Address:  req.User.Address,
		Password: req.NewPassword, // Using NewPassword as the updated password
	}

	res, err := s.service.UpdateUser(r.Context(), req.Id, req.OldPassword, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Using LoginResponse format or similar?
	// The gRPC response had User and Token.
	// We can reuse LoginResponse or a custom struct.
	// Let's reuse LoginResponse for simplicity or define UpdateResponse if needed.
	// But minimal change:
	response := map[string]any{
		"token": res.Token,
		"user":  res.User,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *httpServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// We need password for deletion?
	// Get it from body?
	var req pkg.DeleteReq
	// Ignore error if body is empty? but password is required by service.
	json.NewDecoder(r.Body).Decode(&req)

	if req.Password == "" {
		// try query param
		req.Password = r.URL.Query().Get("password")
	}

	if err := s.service.DeleteUser(r.Context(), id, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (s *httpServer) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := s.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
