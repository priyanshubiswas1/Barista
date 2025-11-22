package pkg

import "time"

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
type RegisterUserWithShopMenuRes struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	ShopId  string `json:"shopId"`
	MenuId  string `json:"menuId"`
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
	ShopId  string `json:"shopId"`
	MenuId  string `json:"menuId"`
}
type UpdateUserReq struct {
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
	Veg         bool    `json:"veg"`
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
	Veg         bool    `json:"veg"`
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
	Veg         bool    `json:"veg"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Available   bool    `json:"available"`
}
type UpdateItemResponse struct {
	Success bool `json:"success"`
}
type UpdateItemAvailableRequest struct {
	ItemId    string `json:"itemId"`
	ShopId    string `json:"shopId"`
	Available bool   `json:"available"`
}
type UpdateItemAvailableResponse struct {
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
type CategoryWithItems struct {
	CategoryId string `json:"categoryId"`
	ShopId     string `json:"shopId"`
	Name       string `json:"name"`
	Items      []Item `json:"items"`
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
type GetCategoriesWithItemsResponse struct {
	Categories []CategoryWithItems `json:"categories"`
}
type DeleteCategoryRequest struct {
	CategoryId string `json:"categoryId"`
	ShopId     string `json:"shopId"`
}
type DeleteCategoryResponse struct {
	Success bool `json:"success"`
}

//IMAGE UPLOAD
type UploadImageResponse struct {
	ImageUrl string `json:"imageUrl"`
}

//TABLES
type Table struct {
	Id          string `json:"tableId"`
	ShopId      string `json:"shopId"`
	TableNumber string `json:"tableNumber"`
	QrUrl       string `json:"qrUrl"`
	Status      string `json:"status"`
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
	PhoneNumber   string  `json:"phoneNumber"`
	Active        bool    `json:"active"`
}
type ActiveSession struct {
	Id            string  `json:"sessionId"`
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	TableNumber   string  `json:"tableNumber"`
	Otp           string  `json:"otp"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	PhoneNumber   string  `json:"phoneNumber"`
	Active        bool    `json:"active"`
}

type GetAllSessionsResponse struct {
	Sessions []Session `json:"sessions"`
}

type GetAllActiveSessionsResponse struct {
	Sessions []ActiveSession `json:"sessions"`
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
	PhoneNumber string `json:"phoneNumber"`
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
type CreateSessionAdminRequest struct {
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	Status        string  `json:"status"`
	Amount        float32 `json:"amount"`
	IsPaid        bool    `json:"isPaid"`
	PaymentMethod string  `json:"paymentMethod"`
	PhoneNumber   string  `json:"phoneNumber"`
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

//ORDERS
type Order struct {
	Id           string `json:"orderId"`
	SessionId    string `json:"sessionId"`
	OrderStatus  string `json:"orderStatus"`
	SpecialNotes string `json:"specialNotes"`
}
type OrderItems struct {
	Id               string  `json:"orderItemId"`
	OrderId          string  `json:"orderId"`
	OrderNumber      string  `json:"orderNumber"`
	ItemId           string  `json:"itemId"`
	Quantity         int32   `json:"quantity"`
	PriceAtOrderTime float32 `json:"priceAtOrderTime"`
	Status           string  `json:"orderStatus"`
}
type OrderItemRequest struct {
	ItemId           string  `json:"itemId"`
	Quantity         int32   `json:"quantity"`
	PriceAtOrderTime float32 `json:"priceAtOrderTime"`
}
type GetOrdersBySessionIdResponse struct {
	OrderItems []OrderItems `json:"orderItems"`
	TotalPrice float32      `json:"totalPrice"`
}
type GetOrdersByOrderIdResponse struct {
	OrderItems []OrderItems `json:"orderItems"`
	TotalPrice float32      `json:"totalPrice"`
}
type CreateOrderAdminRequest struct {
	ShopId       string             `json:"shopId"`
	TableId      string             `json:"tableId"`
	SpecialNotes string             `json:"specialNotes"`
	OrderItems   []OrderItemRequest `json:"orderItems"`
}
type CreateOrderRequest struct {
	SpecialNotes string             `json:"specialNotes"`
	OrderItems   []OrderItemRequest `json:"orderItems"`
}
type CreateOrderResponse struct {
	Order      Order        `json:"order"`
	OrderItems []OrderItems `json:"orderItems"`
}
type UpdateOrderRequest struct {
	OrderId     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"`
}
type UpdateOrderResponse struct {
	Success bool `json:"success"`
}
type CancelOrderRequest struct {
	OrderId string `json:"orderId"`
}
type CancelOrderResponse struct {
	Success bool `json:"success"`
}

//PAYMENTS
type Payment struct {
	Id            string  `json:"id"`
	SessionId     string  `json:"sessionId"`
	Amount        float32 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	PaymentStatus string  `json:"paymentStatus"`
	TransactionId string  `json:"transactionId"`
}
type PaymentRequest struct {
	SessionId     string  `json:"sessionId"`
	ShopId        string  `json:"shopId"`
	TableId       string  `json:"tableId"`
	PhoneNumber   string  `json:"phoneNumber"`
	Amount        float32 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	PaymentStatus string  `json:"paymentStatus"`
	TransactionId string  `json:"transactionId"`
}
type GetPaymentRequest struct {
	SessionId string `json:"sessionId"`
}

//HISTORY
type OrderedItems struct {
	OrderItemId string  `json:"orderItemId"`
	ItemPrice   float32 `json:"itemPrice"`
	Quantity    int32   `json:"quantity"`
}

type SessionSummary struct {
	Id            string         `json:"id"`
	ShopId        string         `json:"shopId"`
	SessionId     string         `json:"sessionId"`
	SessionNumber int32          `json:"sessionNumber"`
	TotalAmount   float32        `json:"totalAmount"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	OrderedItems  []OrderedItems `json:"orderedItems"`
}

type StatusRespone struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetSessionSummaryResponse struct {
	SessionSummaries []SessionSummary `json:"sessionSummaries"`
	TotalAmount      int32            `json:"totalAmount"`
	TotalSessions    int32            `json:"totalSessions"`
}

type CreateSessionSummaryRequest struct {
	SessionId     string `json:"sessionId"`
	ShopId        string `json:"shopId"`
	SessionNumber int32  `json:"sessionNumber"`
}

type CreateSessionSummaryResponse struct {
	SessionSummary SessionSummary `json:"sessionSummary"`
}

type CloseSessionSummaryRequest struct {
	SessionId string `json:"sessionId"`
}

type AddOrderSummaryRequest struct {
	SessionId    string         `json:"sessionId"`
	OrderedItems []OrderedItems `json:"orderedItems"`
}

type RemoveOrderSummaryRequest struct {
	SessionId    string   `json:"sessionId"`
	OrderItemIds []string `json:"orderItemIds"`
}
