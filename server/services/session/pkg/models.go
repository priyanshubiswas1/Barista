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
	PhoneNumber   string  `json:"phoneNumber"`
}

type ActiveSession struct {
	Id            string  `json:"sessionId"`
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	TableName     string  `json:"tableName"`
	Otp           string  `json:"otp"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	Active        bool    `json:"active"`
	PhoneNumber   string  `json:"phoneNumber"`
}

type CreateResponse struct {
	Token string `json:"token"`
}

type Table struct {
	Id          string `json:"tableId"`
	ShopId      string `json:"shopId"`
	TableNumber string `json:"tableNumber"`
	QrUrl       string `json:"qrUrl"`
}

type GetTableResponse struct {
	Table Table `json:"table"`
}
