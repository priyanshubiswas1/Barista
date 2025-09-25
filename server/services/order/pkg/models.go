package pkg

type Order struct {
	Id           string `json:"orderId"`
	SessionId    string `json:"sessionId"`
	OrderStatus  string `json:"orderStatus"`
	SpecialNotes string `json:"specialNotes"`
}

type OrderItems struct {
	Id               string `json:"orderItemId"`
	OrderId          string `json:"orderId"`
	ItemId           string `json:"itemId"`
	Quantity         string `json:"quantity"`
	PriceAtOrderTime string `json:"priceAtOrderTime"`
}

type OrderItemsReq struct {
	ItemId           string `json:"itemId"`
	Quantity         string `json:"quantity"`
	PriceAtOrderTime string `json:"priceAtOrderTime"`
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
