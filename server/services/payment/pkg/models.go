package pkg

type Payment struct {
	Id            string  `json:"id"`
	SessionId     string  `json:"sessionId"`
	Amount        float32 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	Status        string  `json:"status"`
	TransactionId string  `json:"transactionId"`
}

type Session struct {
	Id            string  `json:"sessionId"`
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	Otp           string  `json:"otp"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	PhoneNumber   string  `json:"phoneNumber"`
	Active        bool    `json:"active"`
}

type UpdateSessionRequest struct {
	Session Session `json:"session"`
}
