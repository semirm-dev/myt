package product

import (
	"context"
	pbDiscount "github.com/semirm-dev/myt/discount/proto"
	"github.com/semirm-dev/myt/internal/grpc"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	"github.com/sirupsen/logrus"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "product service"

func NewService(addr string, repo Repository, discountClientAddr string) *defaultService {
	return &defaultService{
		addr:               addr,
		repo:               repo,
		discountClientAddr: discountClientAddr,
	}
}

// defaultService will expose products service via grpc
type defaultService struct {
	pbProduct.UnimplementedProductServer
	addr               string
	repo               Repository
	discountClientAddr string
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

	conn := grpc.CreateClientConnection(svc.discountClientAddr)
	discountsClient := pbDiscount.NewDiscountProviderClient(conn)

	discounts, err := discountsClient.GetDiscounts(ctx, &pbDiscount.DiscountsRequest{})
	if err != nil {
		return nil, err
	}

	logrus.Info(discounts)

	productsResp := toProtoProductsResponse(products)

	return productsResp, nil
}
