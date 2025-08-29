package handler

import (
	menu "github.com/priyanshubiswas1/baristacontracts/gen/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IMenuClient interface {
	// CreateMenu(ctx context.Context,)
}

type MenuClient struct {
	client  *grpc.ClientConn
	service menu.MenuServiceClient
}

func NewMenu(url string) (IMenuClient, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &MenuClient{
		client:  client,
		service: menu.NewMenuServiceClient(client),
	}, nil
}
