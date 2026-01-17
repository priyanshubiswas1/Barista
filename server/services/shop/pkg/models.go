package pkg

type Shop struct {
	ShopId   string `json:"shopId"`
	ShopName string `json:"shopName"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type CreateShopReq struct {
	OwnerId  string `json:"ownerId"`
	ShopName string `json:"shopName"`
}

type UpdateShopReq struct {
	OwnerId string `json:"ownerId"`
	Shop    Shop   `json:"shop"`
}
