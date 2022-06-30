package product

import (
	"context"
	"github.com/semirm-dev/myt/internal/grpc"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "product service"

// defaultService will expose products service via grpc
type defaultService struct {
	pbProduct.UnimplementedProductServer
	addr             string
	repo             Repository
	discountProvider DiscountProvider
}

// Repository communicates to data store with products
type Repository interface {
	GetProductsByFilter(ctx context.Context, filter *Filter) ([]*Product, error)
}

// DiscountProvider communicates to discount service. It is responsible for discounts on products.
type DiscountProvider interface {
	ApplyDiscount(ctx context.Context, products []*pbProduct.ProductMessage) ([]*pbProduct.ProductMessage, error)
}

func NewService(addr string, repo Repository, discountProvider DiscountProvider) *defaultService {
	return &defaultService{
		addr:             addr,
		repo:             repo,
		discountProvider: discountProvider,
	}
}

// ListenForConnections will start grpc server and start listening for connections
func (svc *defaultService) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.addr, serviceName)
}

func (svc *defaultService) RegisterGrpcServer(server *grpcLib.Server) {
	pbProduct.RegisterProductServer(server, svc)
}

// GetProductsByFilter will get filtered products from data store, discounts included
func (svc *defaultService) GetProductsByFilter(ctx context.Context, req *pbProduct.GetProductsByFilterRequest) (*pbProduct.ProductsResponse, error) {
	products, err := svc.repo.GetProductsByFilter(ctx, &Filter{
		PriceLessThan: int(req.PriceLessThen),
		ByCategory:    req.ByCategory,
	})
	if err != nil {
		return nil, err
	}

	productsProto := toProtoProductsMessage(products)

	productsWithDiscount, err := svc.discountProvider.ApplyDiscount(ctx, productsProto)
	if err != nil {
		return nil, err
	}

	return &pbProduct.ProductsResponse{
		Products: productsWithDiscount,
	}, nil
}
