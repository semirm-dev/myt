package product

import (
	pbProduct "github.com/semirm-dev/myt/product/proto"
	"strconv"
)

func toProtoProductMessage(product *Product) *pbProduct.ProductMessage {
	discountPercentage := ""

	if product.Price.DiscountPercentage > 0 {
		discountPercentage = strconv.Itoa(product.Price.DiscountPercentage)
	}

	return &pbProduct.ProductMessage{
		Sku:      product.Sku,
		Name:     product.Name,
		Category: product.Category,
		Price: &pbProduct.PriceMessage{
			Original:           int64(product.Price.Original),
			Final:              int64(product.Price.Final),
			DiscountPercentage: discountPercentage,
			Currency:           product.Price.Currency,
		},
	}
}

func toProtoProductsMessage(products []*Product) []*pbProduct.ProductMessage {
	var protoProducts []*pbProduct.ProductMessage

	for _, product := range products {
		protoProducts = append(protoProducts, toProtoProductMessage(product))
	}

	return protoProducts
}
