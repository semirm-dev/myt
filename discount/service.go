package discount

import (
	"context"
	"fmt"
	pbDiscount "github.com/semirm-dev/myt/discount/proto"
	"github.com/semirm-dev/myt/internal/grpc"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	grpcLib "google.golang.org/grpc"
	"sort"
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
	// we can filter discounts by product.sku or product.category, so we get only discounts applicable to needed products
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
	var filtered []*pbProduct.ProductMessage

	for _, p := range products {
		discountsForProduct := discountsPerProduct(p, discounts)
		distinct := distinctDiscounts(discountsForProduct)

		for _, d := range distinct {
			if d.Percentage > 0 {
				p.Price.DiscountPercentage = fmt.Sprintf("%d%s", d.Percentage, "%")
			}

			if p.Sku == d.Sku {
				p.Price.Final = int64(calculateDiscount(int(p.Price.Original), d.Percentage))
			}

			if p.Category == d.Category {
				p.Price.Final = int64(calculateDiscount(int(p.Price.Original), d.Percentage))
			}
		}

		filtered = append(filtered, p)
	}

	return filtered
}

func distinctDiscounts(discounts []*Discount) []*Discount {
	sort.Slice(discounts, func(i, j int) bool {
		return discounts[i].Percentage > discounts[j].Percentage
	})

	return discounts
}

func discountsPerProduct(product *pbProduct.ProductMessage, discounts []*Discount) []*Discount {
	var filtered []*Discount

	for _, d := range discounts {
		if product.Sku == d.Sku || product.Category == d.Category {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

func calculateDiscount(price int, percentage int) int {
	return price - (price * percentage / 100)
}
