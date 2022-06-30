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
	GetDiscounts(ctx context.Context, filter *Filter) ([]*Discount, error)
}

// Filter discounts from data store, so we do not retrieve all discounts.
// Get only discounts which are related to given products.
type Filter struct {
	Sku      []string
	Category []string
}

func (svc *defaultService) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.addr, serviceName)
}

func (svc *defaultService) RegisterGrpcServer(server *grpcLib.Server) {
	pbDiscount.RegisterDiscountProviderServer(server, svc)
}

// ApplyDiscount will apply discounts on provided products.
func (svc *defaultService) ApplyDiscount(ctx context.Context, req *pbDiscount.DiscountsRequest) (*pbDiscount.DiscountsResponse, error) {
	discounts, err := svc.repo.GetDiscounts(ctx, uniqueFilters(req.Products))
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
		productDiscounts := discountsForProduct(p, discounts)

		// sorting is needed to get the highest discount applied as the final discount
		// aggregate function like MAX() would solve this issue
		sort.Slice(productDiscounts, func(i, j int) bool {
			return productDiscounts[i].Percentage < productDiscounts[j].Percentage
		})

		for _, d := range productDiscounts {
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

		if p.Price.Final == 0 {
			p.Price.Final = p.Price.Original
		}

		filtered = append(filtered, p)
	}

	return filtered
}

// discountsForProduct will filter discounts only related to given product.
// We do not want to apply discount from other products.
// Joins in database would solve this issue.
func discountsForProduct(product *pbProduct.ProductMessage, discounts []*Discount) []*Discount {
	var filtered []*Discount

	for _, d := range discounts {
		if product.Sku == d.Sku || product.Category == d.Category {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

// uniqueFilters will get unique list of product skus and categories.
// This filter is then used to get discounts from data store.
// Where condition in database query would usually solve this issue.
func uniqueFilters(products []*pbProduct.ProductMessage) *Filter {
	filter := &Filter{}
	keys := make(map[string]bool)

	for _, p := range products {
		// sku is unique, no need to check if it's already added to filter
		filter.Sku = append(filter.Sku, p.Sku)

		// filter unique categories from given products
		if _, value := keys[p.Category]; !value {
			keys[p.Category] = true
			filter.Category = append(filter.Category, p.Category)
		}
	}

	return filter
}

func calculateDiscount(price int, percentage int) int {
	return price - (price * percentage / 100)
}
