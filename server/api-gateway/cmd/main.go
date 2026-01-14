package main

import (
	"log"

	"github.com/priyanshubiswas1/barista/apigateway/pkg"
)

func main() {
	cfg, err := pkg.GetConfig()
	if err != nil {
		log.Fatal("failed to get config: ", err)
	}

	// Environment Variable Validation
	validateEnv(cfg)

	log.Println("config loaded")

	server, err := NewAPIServer(cfg.AuthService, cfg.ShopService, cfg.MenuService, cfg.TableService, cfg.SessionService, cfg.OrderService, cfg.PaymentService, cfg.HistoryService, cfg.AllowedOrigins)
	if err != nil {
		log.Fatal("failed to create server")
	}
	log.Println("server created: ", server.authServiceUrl, server.shopServiceUrl, server.menuServiceUrl, server.tableServiceUrl, server.sessionServiceUrl, server.orderServiceUrl, server.paymentServiceUrl, server.historyServiceUrl)

	if err = server.ListenServer(cfg.Port); err != nil {
		log.Fatal("error while listening to port : ", err)
	}
	log.Println("listening to port: ", cfg.Port)
}

func validateEnv(cfg *pkg.Config) {
	required := map[string]string{
		"PORT":                cfg.Port,
		"AUTH_SERVICE_URL":    cfg.AuthService,
		"SHOP_SERVICE_URL":    cfg.ShopService,
		"MENU_SERVICE_URL":    cfg.MenuService,
		"TABLE_SERVICE_URL":   cfg.TableService,
		"SESSION_SERVICE_URL": cfg.SessionService,
		"ORDER_SERVICE_URL":   cfg.OrderService,
		"PAYMENT_SERVICE_URL": cfg.PaymentService,
		"HISTORY_SERVICE_URL": cfg.HistoryService,
		"JWT_KEY":             cfg.JwtKey,
	}

	for name, val := range required {
		if val == "" {
			log.Fatalf("Environment variable %s is required", name)
		}
	}
}
