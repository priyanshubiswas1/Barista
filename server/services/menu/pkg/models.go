package pkg

type Menu struct {
	MenuId string `json:"menuId"`
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}

type ItemRequest struct {
	ItemId      string  `json:"menuItemId"`
	ShopId      string  `json:"shopId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryId  string  `json:"categoryId"`
	Veg         bool    `json:"veg"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	IsAvailable bool    `json:"isAvailable"`
}

type Item struct {
	ItemId      string  `json:"menuItemId"`
	ShopId      string  `json:"shopId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Veg         bool    `json:"veg"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	IsAvailable bool    `json:"isAvailable"`
}

type MenuItem struct {
	Id     string `json:"id"`
	MenuId string `json:"menuId"`
	ItemId string `json:"ItemId"`
}

type Category struct {
	Id     string `json:"id"`
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}

type CategoryWithItems struct {
	Id     string `json:"id"`
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
	Items  []Item `json:"items"`
}
