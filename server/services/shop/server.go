package shop

import (
	"context"
	"fmt"
	"net"

	"net/http"

	"github.com/priyanshubiswas1/barista/shopservice/pkg"
	shopPb "github.com/priyanshubiswas1/baristacontracts/gen/shop"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	shopPb.UnimplementedShopServiceServer
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
	shopPb.RegisterShopServiceServer(grpcS, &gRPCserver{service: service})
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

	fmt.Printf("shop service listening on %s (gRPC and HTTP health checks)\n", port)
	return m.Serve()
}

func (s *gRPCserver) GetShopByOwnerId(ctx context.Context, req *shopPb.GetByOwnerIdRequest) (*shopPb.GetByOwnerIdResponse, error) {
	shops, err := s.service.GetByOwnerId(ctx, req.OwnerId)
	if err != nil {
		return nil, err
	}
	protoShops := []*shopPb.Shop{}
	for _, s := range shops {
		protoShops = append(protoShops, &shopPb.Shop{
			ShopId:   s.ShopId,
			ShopName: s.ShopName,
		})
	}
	return &shopPb.GetByOwnerIdResponse{
		Shops: protoShops,
	}, nil
}

func (s *gRPCserver) GetShopById(ctx context.Context, req *shopPb.GetByShopIdRequest) (*shopPb.GetShopByIdResponse, error) {
	res, err := s.service.GetByShopId(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	return &shopPb.GetShopByIdResponse{
		Shop: &shopPb.Shop{
			ShopId:   res.ShopId,
			ShopName: res.ShopName,
		},
	}, nil
}

func (s *gRPCserver) RegisterShop(ctx context.Context, req *shopPb.CreateRequest) (*shopPb.CreateResponse, error) {
	id, err := s.service.CreateShop(ctx, req.OwnerId, req.ShopName)
	if err != nil {
		return nil, err
	}
	return &shopPb.CreateResponse{
		Shop: &shopPb.Shop{
			ShopId:   id,
			ShopName: req.ShopName,
		},
	}, nil
}
func (s *gRPCserver) UpdateShop(ctx context.Context, req *shopPb.UpdateRequest) (*shopPb.UpdateResponse, error) {
	if err := s.service.UpdateShop(ctx, req.OwnerId, &pkg.Shop{
		ShopId:   req.Shop.ShopId,
		ShopName: req.Shop.ShopName,
	}); err != nil {
		return nil, err
	}
	return &shopPb.UpdateResponse{
		Shop: &shopPb.Shop{
			ShopId:   req.Shop.ShopId,
			ShopName: req.Shop.ShopName,
		},
	}, nil
}
func (s *gRPCserver) DeleteShop(ctx context.Context, req *shopPb.DeleteRequest) (*shopPb.DeleteResponse, error) {
	if err := s.service.DeleteShop(ctx, req.OwnerId, req.ShopId); err != nil {
		return &shopPb.DeleteResponse{
			Success: false,
		}, err
	}
	return &shopPb.DeleteResponse{
		Success: true,
	}, nil
}
