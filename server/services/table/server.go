package table

import (
	"context"
	"fmt"
	"net"

	"net/http"

	"github.com/priyanshubiswas1/barista/tableservice/pkg"
	tablePb "github.com/priyanshubiswas1/baristacontracts/gen/table"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	tablePb.UnimplementedTableServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}

	m := cmux.New(lis)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := grpc.NewServer()
	tablePb.RegisterTableServiceServer(grpcS, &gRPCserver{
		service: service,
	})
	reflection.Register(grpcS)

	httpS := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
			if r.URL.Path == "/ready" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Ready"))
				return
			}
			http.NotFound(w, r)
		}),
	}

	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	fmt.Printf("table service listening on %s (gRPC and HTTP health checks)\n", port)
	return m.Serve()
}

func (s *gRPCserver) CreateTable(ctx context.Context, req *tablePb.CreateTableRequest) (*tablePb.CreateTableResponse, error) {
	pkgTable, err := s.service.CreateTable(ctx, req.ShopId, req.TableNumber)
	if err != nil {
		return nil, err
	}
	return &tablePb.CreateTableResponse{
		Table: &tablePb.Table{
			Id:          pkgTable.TableId,
			ShopId:      pkgTable.ShopId,
			TableNumber: pkgTable.TableNumber,
			Status:      pkgTable.Status,
			Qr:          pkgTable.QR,
		},
	}, nil
}

func (s *gRPCserver) GetTable(ctx context.Context, req *tablePb.GetTableRequest) (*tablePb.GetTableResponse, error) {
	pkgTable, err := s.service.GetTable(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &tablePb.GetTableResponse{
		Table: &tablePb.Table{
			Id:          pkgTable.TableId,
			ShopId:      pkgTable.ShopId,
			TableNumber: pkgTable.TableNumber,
			Status:      pkgTable.Status,
			Qr:          pkgTable.QR,
		},
	}, nil
}

func (s *gRPCserver) GetTables(ctx context.Context, req *tablePb.GetTablesRequest) (*tablePb.GetTablesResponse, error) {
	pkgTables, err := s.service.GetTables(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	tables := []*tablePb.Table{}
	for _, pkgTable := range pkgTables {
		tables = append(tables, &tablePb.Table{
			Id:          pkgTable.TableId,
			ShopId:      pkgTable.ShopId,
			TableNumber: pkgTable.TableNumber,
			Status:      pkgTable.Status,
			Qr:          pkgTable.QR,
		})
	}
	return &tablePb.GetTablesResponse{
		Tables: tables,
	}, nil
}

func (s *gRPCserver) UpdateTable(ctx context.Context, req *tablePb.UpdateTableRequest) (*tablePb.UpdateTableResponse, error) {
	if err := s.service.UpdateTable(ctx, pkg.Table{
		TableId:     req.Table.Id,
		ShopId:      req.Table.ShopId,
		TableNumber: req.Table.TableNumber,
		Status:      req.Table.Status,
		QR:          req.Table.Qr,
	}); err != nil {
		return nil, err
	}
	return &tablePb.UpdateTableResponse{
		Success: true,
	}, nil
}

func (s *gRPCserver) UpdateTableStatus(ctx context.Context, req *tablePb.UpdateTableStatusRequest) (*tablePb.UpdateTableStatusResponse, error) {
	if err := s.service.UpdateTableStatus(ctx, pkg.Table{
		TableId: req.TableId,
		ShopId:  req.ShopId,
		Status:  req.Status,
	}); err != nil {
		return nil, err
	}
	return &tablePb.UpdateTableStatusResponse{
		Success: true,
	}, nil
}

func (s *gRPCserver) DeleteTable(ctx context.Context, req *tablePb.DeleteTableRequest) (*tablePb.DeleteTableResponse, error) {
	if err := s.service.DeleteTable(ctx, req.Id, req.ShopId); err != nil {
		return nil, err
	}
	return &tablePb.DeleteTableResponse{
		Success: true,
	}, nil
}
