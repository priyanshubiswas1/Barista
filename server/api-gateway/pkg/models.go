package pkg

//USER
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

//SHOP
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

//MENU
type Menu struct {
	MenuId string `json:"menuId"`
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}
type Item struct {
	ItemId      string  `json:"itemId"`
	ShopId      string  `json:"shopId"`
	ItemName    string  `json:"itemName"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Available   bool    `json:"available"`
}
type MenuItem struct {
	MenuItemId string `json:"menuItemId"`
	MenuId     string `json:"menuId"`
	ItemId     string `json:"itemId"`
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
type SetMenuResponse struct {
	Success bool `json:"success"`
}
type DeleteMenuResponse struct {
	Success bool `json:"success"`
}
type CreateItemRequest struct {
	ItemName    string  `json:"name"`
	ShopId      string  `json:"shopId"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Available   bool    `json:"available"`
}
type CreateItemResponse struct {
	Item Item `json:"item"`
}
type GetItemResponse struct {
	Item Item `json:"item"`
}
type GetItemsResponse struct {
	Items []Item `json:"items"`
}
type UpdateItemRequest struct {
	Item Item `json:"item"`
}
type UpdateItemResponse struct {
	Success bool `json:"success"`
}
type DeleteItemResponse struct {
	Success bool `json:"success"`
}
type CreateMenuItemRequest struct {
	MenuId string `json:"menuId"`
	ItemId string `json:"itemId"`
}
type CreateMenuItemResponse struct {
	MenuItem MenuItem `json:"menuItem"`
}
type GetMenuItemsResponse struct {
	Items []Item `json:"items"`
}
type DeleteMenuItemRequest struct {
	MenuId string `json:"menuId"`
	ItemId string `json:"itemId"`
}
type DeleteMenuItemResponse struct {
	Success bool `json:"success"`
}
type DeleteMenuItemsResponse struct {
	Success bool `json:"success"`
}
