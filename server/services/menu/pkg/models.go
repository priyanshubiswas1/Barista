package pkg

type Menu struct {
	MenuId string `json:"menuId"`
	ShopId string `json:"shopId"`
	Name   string `json:"name"`
}

type MenuItem struct {
	MenuItemId  string  `json:"menuItemId"`
	MenuId      string  `json:"menuId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	IsAvailable bool    `json:"isAvailable"`
}
