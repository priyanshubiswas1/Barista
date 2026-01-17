package pkg

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type RegisterWithShopMenuResponse struct {
	Token  string `json:"token"`
	ShopId string `json:"shopId"`
	MenuId string `json:"menuId"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	User   User   `json:"user"`
	ShopId string `json:"shopId"`
	MenuId string `json:"menuId"`
}

type CreateShopReq struct {
	OwnerId  string `json:"ownerId"`
	ShopName string `json:"shopName"`
}

type CreateMenuReq struct {
	ShopId   string `json:"shopId"`
	MenuName string `json:"menuName"`
}

type RegisterReq struct {
	User     User   `json:"user"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Password    string `json:"password"`
}

type UpdateReq struct {
	Id          string `json:"id"`
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
	User        User   `json:"user"`
}

type DeleteReq struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
