package pkg

type Table struct {
	TableId     string `json:"tableId"`
	ShopId      string `json:"shopId"`
	TableNumber string `json:"tableNumber"`
	Status      string `json:"status"`
	QR          string `json:"qrUrl"`
}
