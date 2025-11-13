package pkg

type Order struct {
	Id           string `json:"orderId"`
	SessionId    string `json:"sessionId"`
	OrderStatus  string `json:"orderStatus"`
	SpecialNotes string `json:"specialNotes"`
}

type OrderItems struct {
	Id               string  `json:"orderItemId"`
	OrderId          string  `json:"orderId"`
	ItemId           string  `json:"itemId"`
	Quantity         int32   `json:"quantity"`
	PriceAtOrderTime float32 `json:"priceAtOrderTime"`
	OrderStatus      string  `json:"orderStatus"`
}

type OrderItemsReq struct {
	ItemId           string  `json:"itemId"`
	Quantity         int32   `json:"quantity"`
	PriceAtOrderTime float32 `json:"priceAtOrderTime"`
}

type OrderAdminReq struct {
	ShopId       string          `json:"shopId"`
	TableId      string          `json:"tableId"`
	SpecialNotes string          `json:"specialNotes"`
	OrderItems   []OrderItemsReq `json:"orderItems"`
}

type OrderReq struct {
	SessionId    string          `json:"sessionId"`
	SpecialNotes string          `json:"specialNotes"`
	OrderItems   []OrderItemsReq `json:"orderItems"`
}

type CreateOrderRes struct {
	Order      Order        `json:"order"`
	OrderItems []OrderItems `json:"orderItems"`
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

type AddOrderHistoryReq struct {
	OrderItemId string  `json:"orderItemId"`
	OrderId     string  `json:"orderId"`
	ItemPrice   float32 `json:"itemPrice"`
	Quantity    int32   `json:"quantity"`
}
