package shop

import (
	"context"
	"fmt"
	"net"

	"github.com/priyanshubiswas1/barista/shopservice/pkg"
	shop "github.com/priyanshubiswas1/baristacontracts/gen/shop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCserver struct {
	service IService
	shop.UnimplementedShopServiceServer
}

func ListenGRPC(service IService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	shop.RegisterShopServiceServer(server, &gRPCserver{service: service})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *gRPCserver) GetShopByOwnerId(ctx context.Context, req *shop.GetByOwnerIdRequest) (*shop.GetByOwnerIdResponse, error) {
	shops, err := s.service.GetByOwnerId(ctx, req.OwnerId)
	if err != nil {
		return nil, err
	}
	protoShops := []*shop.Shop{}
	for _, s := range shops {
		protoShops = append(protoShops, &shop.Shop{
			ShopId:   s.ShopId,
			ShopName: s.ShopName,
		})
	}
	return &shop.GetByOwnerIdResponse{
		Shops: protoShops,
	}, nil
}

func (s *gRPCserver) GetShopById(ctx context.Context, req *shop.GetByShopIdRequest) (*shop.GetShopByIdResponse, error) {
	res, err := s.service.GetByShopId(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}
	return &shop.GetShopByIdResponse{
		Shop: &shop.Shop{
			ShopId:   res.ShopId,
			ShopName: res.ShopName,
		},
	}, nil
}

func (s *gRPCserver) RegisterShop(ctx context.Context, req *shop.CreateRequest) (*shop.CreateResponse, error) {
	id, err := s.service.CreateShop(ctx, req.OwnerId, req.ShopName)
	if err != nil {
		return nil, err
	}
	return &shop.CreateResponse{
		Shop: &shop.Shop{
			ShopId:   id,
			ShopName: req.ShopName,
		},
	}, nil
}
func (s *gRPCserver) UpdateShop(ctx context.Context, req *shop.UpdateRequest) (*shop.UpdateResponse, error) {
	if err := s.service.UpdateShop(ctx, req.OwnerId, &pkg.Shop{
		ShopId:   req.Shop.ShopId,
		ShopName: req.Shop.ShopName,
	}); err != nil {
		return nil, err
	}
	return &shop.UpdateResponse{
		Shop: &shop.Shop{
			ShopId:   req.Shop.ShopId,
			ShopName: req.Shop.ShopName,
		},
	}, nil
}
func (s *gRPCserver) DeleteShop(ctx context.Context, req *shop.DeleteRequest) (*shop.DeleteResponse, error) {
	if err := s.service.DeleteShop(ctx, req.OwnerId, req.ShopId); err != nil {
		return &shop.DeleteResponse{
			Success: false,
		}, err
	}
	return &shop.DeleteResponse{
		Success: true,
	}, nil
}
