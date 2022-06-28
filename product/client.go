package product

import (
	"context"
	"github.com/semirm-dev/myt/internal/grpc"
	pbProduct "github.com/semirm-dev/myt/product/proto"
)

// Client is responsible for communication with product Service
type Client interface {
	GetProductsByFilter(ctx context.Context, filter *Filter) ([]*pbProduct.ProductMessage, error)
}

type Filter struct {
	PriceLessThan int
	ByCategory    string
}

type defaultClient struct {
	grpcClient pbProduct.ProductClient
}

func NewClient(addr string) Client {
	conn := grpc.CreateClientConnection(addr)

	return &defaultClient{
		grpcClient: pbProduct.NewProductClient(conn),
	}
}

func (client *defaultClient) GetProductsByFilter(ctx context.Context, filter *Filter) ([]*pbProduct.ProductMessage, error) {
	productsResp, err := client.grpcClient.GetProductsByFilter(ctx, &pbProduct.GetProductsByFilterRequest{
		PriceLessThen: int64(filter.PriceLessThan),
		ByCategory:    filter.ByCategory,
	})
	if err != nil {
		return nil, err
	}

	return productsResp.Products, nil
}
