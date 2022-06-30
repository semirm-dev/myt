package discount

import (
	"context"
	pbDiscount "github.com/semirm-dev/myt/discount/proto"
	"github.com/semirm-dev/myt/internal/grpc"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "discount service"

func NewService(addr string, repo Repository) *defaultService {
	return &defaultService{
		addr: addr,
		repo: repo,
	}
}

// defaultService will expose discount service via grpc
type defaultService struct {
	pbDiscount.UnimplementedDiscountProviderServer
	addr string
	repo Repository
}

// Repository communicates to data store with discounts
type Repository interface {
	GetDiscount(ctx context.Context) ([]*Discount, error)
}

func (svc *defaultService) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.addr, serviceName)
}

func (svc *defaultService) RegisterGrpcServer(server *grpcLib.Server) {
	pbDiscount.RegisterDiscountProviderServer(server, svc)
}

// ApplyDiscount will apply discounts on provided products.
func (svc *defaultService) ApplyDiscount(ctx context.Context, req *pbDiscount.DiscountsRequest) (*pbDiscount.DiscountsResponse, error) {
	discounts, err := svc.repo.GetDiscount(ctx)
	if err != nil {
		return nil, err
	}

	productsWithDiscount := applyDiscount(req.Products, discounts)

	return &pbDiscount.DiscountsResponse{
		Products: productsWithDiscount,
	}, nil
}

func applyDiscount(products []*pbProduct.ProductMessage, discounts []*Discount) []*pbProduct.ProductMessage {
	return products
}
