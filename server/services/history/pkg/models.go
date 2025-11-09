package pkg

import "time"

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
