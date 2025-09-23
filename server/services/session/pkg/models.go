package pkg

type Session struct {
	Id            string  `json:"sessionId"`
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	Otp           string  `json:"otp"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	Active        bool    `json:"active"`
	PhoneNumber   int32   `json:"phoneNumber"`
}

type CreateResponse struct {
	Token string `json:"token"`
}
