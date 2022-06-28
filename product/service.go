package product

import (
	"context"
	"github.com/semirm-dev/myt/internal/grpc"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "product service"

func NewService(addr string, repo Repository) *defaultService {
	return &defaultService{
		addr: addr,
		repo: repo,
	}
}

// defaultService will expose products service via grpc
type defaultService struct {
	pbProduct.UnimplementedProductServer
	addr string
	repo Repository
}

func (svc *defaultService) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.addr, serviceName)
}

func (svc *defaultService) RegisterGrpcServer(server *grpcLib.Server) {
	pbProduct.RegisterProductServer(server, svc)
}

func (svc *defaultService) GetProductsByFilter(ctx context.Context, req *pbProduct.GetProductsByFilterRequest) (*pbProduct.ProductsResponse, error) {
	products, err := svc.repo.GetProductsByFilter(ctx, &Filter{
		PriceLessThan: int(req.PriceLessThen),
		ByCategory:    req.ByCategory,
	})
	if err != nil {
		return nil, err
	}

	return toProtoProductsResponse(products), nil
}
