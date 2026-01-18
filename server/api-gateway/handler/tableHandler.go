package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priyanshubiswas1/barista/apigateway/pkg"
	table "github.com/priyanshubiswas1/baristacontracts/gen/table"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ITableClient interface {
	CreateTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	GetTables(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	UpdateTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	DeleteTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type TableClient struct {
	client  *grpc.ClientConn
	service table.TableServiceClient
}

func NewTable(url string) (ITableClient, error) {
	client, err := grpc.NewClient(pkg.StripScheme(url), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &TableClient{}, err
	}
	return &TableClient{
		client:  client,
		service: table.NewTableServiceClient(client),
	}, nil
}

func (c *TableClient) CreateTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.CreateTableRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.CreateTable(ctx, &table.CreateTableRequest{
		ShopId:      req.ShopId,
		TableNumber: req.TableNumber,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.CreateTableResponse{
		Table: pkg.Table{
			Id:          res.Table.Id,
			ShopId:      res.Table.ShopId,
			TableNumber: res.Table.TableNumber,
			Status:      res.Table.Status,
			QrUrl:       res.Table.Qr,
		},
	})
}

func (c *TableClient) GetTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	tableId := vars["tableid"]
	res, err := c.service.GetTable(ctx, &table.GetTableRequest{
		Id: tableId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetTableResponse{
		Table: pkg.Table{
			Id:          res.Table.Id,
			ShopId:      res.Table.ShopId,
			TableNumber: res.Table.TableNumber,
			Status:      res.Table.Status,
			QrUrl:       res.Table.Qr,
		},
	})
}

func (c *TableClient) GetTables(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	shopId := vars["shopid"]
	res, err := c.service.GetTables(ctx, &table.GetTablesRequest{
		ShopId: shopId,
	})
	if err != nil {
		return err
	}
	tables := []pkg.Table{}
	for _, t := range res.Tables {
		tables = append(tables, pkg.Table{
			Id:          t.Id,
			ShopId:      t.ShopId,
			TableNumber: t.TableNumber,
			Status:      t.Status,
			QrUrl:       t.Qr,
		})
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.GetTablesResponse{
		Tables: tables,
	})
}

func (c *TableClient) UpdateTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.UpdateTableRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.UpdateTable(ctx, &table.UpdateTableRequest{
		Table: &table.Table{
			Id:          req.Table.Id,
			ShopId:      req.Table.ShopId,
			TableNumber: req.Table.TableNumber,
			Qr:          req.Table.QrUrl,
		},
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.UpdateTableResponse{
		Success: res.Success,
	})
}

func (c *TableClient) DeleteTable(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	req := pkg.DeleteTableRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	res, err := c.service.DeleteTable(ctx, &table.DeleteTableRequest{
		Id:     req.Id,
		ShopId: req.ShopId,
	})
	if err != nil {
		return err
	}
	return pkg.WriteJSON(w, http.StatusOK, pkg.DeleteTableResponse{
		Success: res.Success,
	})
}
