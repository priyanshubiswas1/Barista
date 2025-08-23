package pkg

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
