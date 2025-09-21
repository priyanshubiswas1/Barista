package table

import (
	"context"
	"fmt"
	"net"

	"github.com/priyanshubiswas1/barista/tableservice/pkg"
	"github.com/priyanshubiswas1/baristacontracts/gen/table"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	table.UnimplementedTableServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	table.RegisterTableServiceServer(server, &gRPCserver{
		service: service,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) CreateTable(ctx context.Context, req *table.CreateTableRequest) (*table.CreateTableResponse, error) {
	pkgTable, err := s.service.CreateTable(ctx, req.ShopId, req.TableNumber)
	if err != nil {
		return nil, err
	}
	return &table.CreateTableResponse{
		Table: &table.Table{
			Id:          pkgTable.TableId,
			ShopId:      pkgTable.ShopId,
			TableNumber: pkgTable.TableNumber,
			QrUrl:       pkgTable.QR,
		},
	}, nil
}

func (s *gRPCserver) GetTable(ctx context.Context, req *table.GetTableRequest) (*table.GetTableResponse, error) {
	pkgTable, err := s.service.GetTable(ctx, req.Id, req.ShopId)
	if err != nil {
		return nil, err
	}
	return &table.GetTableResponse{
		Table: &table.Table{
			Id:          pkgTable.TableId,
			ShopId:      pkgTable.ShopId,
			TableNumber: pkgTable.TableNumber,
			QrUrl:       pkgTable.QR,
		},
	}, nil
}

func (s *gRPCserver) GetTables(ctx context.Context, req *table.GetTablesRequest) (*table.GetTablesResponse, error) {
	pkgTables, err := s.service.GetTables(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	tables := []*table.Table{}
	for _, pkgTable := range pkgTables {
		tables = append(tables, &table.Table{
			Id:          pkgTable.TableId,
			ShopId:      pkgTable.ShopId,
			TableNumber: pkgTable.TableNumber,
			QrUrl:       pkgTable.QR,
		})
	}
	return &table.GetTablesResponse{
		Tables: tables,
	}, nil
}

func (s *gRPCserver) UpdateTable(ctx context.Context, req *table.UpdateTableRequest) (*table.UpdateTableResponse, error) {
	if err := s.service.UpdateTable(ctx, pkg.Table{
		TableId:     req.Table.Id,
		ShopId:      req.Table.ShopId,
		TableNumber: req.Table.TableNumber,
		QR:          req.Table.QrUrl,
	}); err != nil {
		return nil, err
	}
	return &table.UpdateTableResponse{
		Success: true,
	}, nil
}

func (s *gRPCserver) DeleteTable(ctx context.Context, req *table.DeleteTableRequest) (*table.DeleteTableResponse, error) {
	if err := s.service.DeleteTable(ctx, req.Id, req.ShopId); err != nil {
		return nil, err
	}
	return &table.DeleteTableResponse{
		Success: true,
	}, nil
}
