package client

import (
	"context"

	"errors"

	"github.com/priyanshubiswas1/barista/sessionservice/pkg"
	table "github.com/priyanshubiswas1/baristacontracts/gen/table"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ITableClient interface {
	GetTable(ctx context.Context, tableId, shopId string) (pkg.GetTableResponse, error)
	UpdateTableStatus(ctx context.Context, tableId, shopId, status string) (pkg.UpdateTableStatusResponse, error)
}

type TableClient struct {
	client  *grpc.ClientConn
	service table.TableServiceClient
}

func NewTable(url string) (ITableClient, error) {
	if url == "" {
		return &TableClient{}, nil
	}
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &TableClient{
		client:  client,
		service: table.NewTableServiceClient(client),
	}, nil
}

func (c *TableClient) GetTable(ctx context.Context, tableId, shopId string) (pkg.GetTableResponse, error) {
	if c.service == nil {
		return pkg.GetTableResponse{}, errors.New("table service URL not configured")
	}
	res, err := c.service.GetTable(ctx, &table.GetTableRequest{
		Id:     tableId,
		ShopId: shopId,
	})
	if err != nil {
		return pkg.GetTableResponse{}, err
	}
	return pkg.GetTableResponse{
		Table: pkg.Table{
			Id:          res.Table.Id,
			ShopId:      res.Table.ShopId,
			TableNumber: res.Table.TableNumber,
			QrUrl:       res.Table.Qr,
		},
	}, nil
}

func (c *TableClient) UpdateTableStatus(ctx context.Context, tableId, shopId, status string) (pkg.UpdateTableStatusResponse, error) {
	if c.service == nil {
		return pkg.UpdateTableStatusResponse{}, errors.New("table service URL not configured")
	}
	res, err := c.service.UpdateTableStatus(ctx, &table.UpdateTableStatusRequest{
		TableId: tableId,
		ShopId:  shopId,
		Status:  status,
	})
	if err != nil {
		return pkg.UpdateTableStatusResponse{}, err
	}
	return pkg.UpdateTableStatusResponse{
		Success: res.Success,
	}, nil
}
