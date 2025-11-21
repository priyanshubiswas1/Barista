package pkg

import "time"

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

type UpdateTableStatusResponse struct {
	Success bool `json:"success"`
}

type OrderedItem struct {
	OrderItemId string  `json:"orderItemId"`
	OrderId     string  `json:"orderId"`
	ItemPrice   float32 `json:"itemPrice"`
	Quantity    int32   `json:"quantity"`
}

type Summary struct {
	Id            string        `json:"summaryId"`
	ShopId        string        `json:"shopId"`
	SessionId     string        `json:"sessionId"`
	SessionNumber int32         `json:"sessionNumber"`
	TotalAmount   float32       `json:"totalAmount"`
	OrderedItems  []OrderedItem `json:"orderedItems"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}
