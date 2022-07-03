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
	discounts, err := svc.repo.GetDiscounts(ctx, filterDiscountsFor(req.Products))
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
		// instead of individually selecting discounts for each product in for loop, 
		// we could also have these grouped in previous steps (select from data base)
		/*
			"product-sku-1": [
				{
					"discount1"
				},
				{
					"discount2"
				}
			]
		*/
		productDiscounts := discountsPerProduct(p, discounts)

		// sorting is needed to get the highest discount applied as the final discount
		sort.Slice(productDiscounts, func(i, j int) bool {
			return productDiscounts[i].Percentage > productDiscounts[j].Percentage
		})

		if len(productDiscounts) > 0 {
			biggestDiscount := productDiscounts[0]

			if biggestDiscount.Percentage > 0 {
				p.Price.DiscountPercentage = fmt.Sprintf("%d%s", biggestDiscount.Percentage, "%")
			}

			if p.Category == biggestDiscount.Category {
				p.Price.Final = int64(calculateDiscount(int(p.Price.Original), biggestDiscount.Percentage))
			}

			if p.Sku == biggestDiscount.Sku {
				p.Price.Final = int64(calculateDiscount(int(p.Price.Original), biggestDiscount.Percentage))
			}
		}

		if p.Price.Final == 0 {
			p.Price.Final = p.Price.Original
		}

		filtered = append(filtered, p)
	}

	return filtered
}

// discountsPerProduct will filter discounts only related to given product.
// We do not want to apply discount from other products.
func discountsPerProduct(product *pbProduct.ProductMessage, discounts []*Discount) []*Discount {
	var filtered []*Discount

	for _, d := range discounts {
		if product.Sku == d.Sku || product.Category == d.Category {
			filtered = append(filtered, d)
		}
	}

	return filtered
}

// filterDiscountsFor will get unique list of products' skus and categories.
// This filter is then used to get discounts from data store.
func filterDiscountsFor(products []*pbProduct.ProductMessage) *Filter {
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
