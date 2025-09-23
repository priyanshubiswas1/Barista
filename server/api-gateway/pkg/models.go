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
	CategoryId  string  `json:"categoryId"`
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
	ItemId      string  `json:"itemId"`
	ShopId      string  `json:"shopId"`
	ItemName    string  `json:"itemName"`
	Description string  `json:"description"`
	CategoryId  string  `json:"categoryId"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Available   bool    `json:"available"`
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

//CATEGORY
type Category struct {
	CategoryId string `json:"categoryId"`
	ShopId     string `json:"shopId"`
	Name       string `json:"name"`
}
type CreateCategoryRequest struct {
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}
type CreateCategoryResponse struct {
	Category Category `json:"category"`
}
type GetCategoriesRequest struct {
	ShopId string `json:"shopId"`
}
type GetCategoriesResponse struct {
	Categories []Category `json:"categories"`
}
type DeleteCategoryRequest struct {
	CategoryId string `json:"categoryId"`
	ShopId     string `json:"shopId"`
}
type DeleteCategoryResponse struct {
	Success bool `json:"success"`
}

//TABLES
type Table struct {
	Id          string `json:"tableId"`
	ShopId      string `json:"shopId"`
	TableNumber string `json:"tableNumber"`
	QrUrl       string `json:"qrUrl"`
}
type CreateTableRequest struct {
	ShopId      string `json:"shopId"`
	TableNumber string `json:"tableNumber"`
}
type CreateTableResponse struct {
	Table Table `json:"table"`
}
type GetTablesRequest struct {
	ShopId string `json:"shopId"`
}
type GetTablesResponse struct {
	Tables []Table `json:"tables"`
}
type GetTableRequest struct {
	Id     string `json:"tableId"`
	ShopId string `json:"shopId"`
}
type GetTableResponse struct {
	Table Table `json:"table"`
}
type UpdateTableRequest struct {
	Table Table `json:"table"`
}
type UpdateTableResponse struct {
	Success bool `json:"success"`
}
type DeleteTableRequest struct {
	Id     string `json:"tableId"`
	ShopId string `json:"shopId"`
}
type DeleteTableResponse struct {
	Success bool `json:"success"`
}

//SESSIONS
type Session struct {
	Id            string  `json:"sessionId"`
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	Otp           string  `json:"otp"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	PhoneNumber   int32   `json:"phoneNumber"`
	Active        bool    `json:"active"`
}

type GetAllSessionsResponse struct {
	Sessions []Session `json:"sessions"`
}
type GetActiveSessionResponse struct {
	Session Session `json:"session"`
}
type GetSessionRequest struct {
	Id string `json:"sessionId"`
}
type GetSessionResponse struct {
	Session Session `json:"session"`
}
type SendSessionOTPRequest struct {
	ShopId      string `json:"shopId"`
	TableId     string `json:"tableId"`
	PhoneNumber int32  `json:"phoneNumber"`
}
type SendSessionOTPResponse struct {
	Success bool `json:"success"`
}
type CreateSessionRequest struct {
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	Otp           string  `json:"otp"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	Active        bool    `json:"active"`
}
type CreateSessionResponse struct {
	Session      Session `json:"session"`
	SessionToken string  `json:"sessionToken"`
}
type JoinSessionRequest struct {
	ShopId  string `json:"shopId"`
	TableId string `json:"tableId"`
	Otp     string `json:"otp"`
}
type JoinSessionResponse struct {
	Session      Session `json:"session"`
	SessionToken string  `json:"sessionToken"`
}
type UpdateSessionRequest struct {
	Session Session `json:"session"`
}
type UpdateSessionResponse struct {
	Success bool `json:"success"`
}
