package pkg

type User struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type RegisterRes struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UpdateUserReq struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type UpdateUserRes struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type DeleteUserReq struct {
	Id string `json:"id"`
}

type DeleteUserRes struct {
	Success bool `json:"success"`
}

//////////////////////////////////////////////////////////////////////////

type Shop struct {
	ShopName string `json:"shopName"`
	ShopId   string `json:"shopId"`
}

type GetByShopIdRequest struct {
	ShopId string `json:"shopId"`
}

type CreateShopReq struct {
	ShopName string `json:"shopName"`
}

type CreateShopRes struct {
	Shop Shop `json:"shop"`
}

type UpdateShopReq struct {
	Shop Shop `json:"shop"`
}

type UpdateShopRes struct {
	Shop Shop `json:"shop"`
}

type DeleteShopReq struct {
	ShopId string `json:"shopId"`
}

type DeleteShopRes struct {
	Success bool `json:"success"`
}
