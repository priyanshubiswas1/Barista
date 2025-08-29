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
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
	Address     string `json:"address"`
}

type UpdateUserRes struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type DeleteUserReq struct {
	Password string `json:"password"`
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

//////////////////////////////////////////////////////////////////////////

type Menu struct {
	MenuId string `json:"menuId"`
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}

type MenuItem struct {
	MenuItemId  string  `json:"menuItemId"`
	MenuId      string  `json:"menuId"`
	ItemName    string  `json:"itemName"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Available   bool    `json:"available"`
}

type CreateMenuRequest struct {
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}

type CreateMenuResponse struct {
	Menu Menu `json:"menu"`
}

type GetMenusResponse struct {
	Menus []Menu `json:"menus"`
}

type DeleteMenuResponse struct {
	Success bool `json:"success"`
}

type CreateMenuItemRequest struct {
	ItemName    string  `json:"itemName"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Available   bool    `json:"available"`
}

type CreateMenuItemResponse struct {
	MenuItem MenuItem `json:"menuItem"`
}

type GetMenuItemsResponse struct {
	MenuItems []MenuItem `json:"menuItems"`
}

type GetMenuItemResponse struct {
	MenuItem MenuItem `json:"menuItem"`
}

type UpdateMenuItemRequest struct {
	MenuItem MenuItem `json:"menuItem"`
}

type UpdateMenuItemResponse struct {
	Success bool `json:"success"`
}

type DeleteMenuItemsResponse struct {
	Success bool `json:"success"`
}

type DeleteMenuItemResponse struct {
	Success bool `json:"success"`
}
