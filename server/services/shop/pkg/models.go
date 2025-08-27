package pkg

type Shop struct {
	ShopId   string `json:"shopId"`
	ShopName string `json:"shopName"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
