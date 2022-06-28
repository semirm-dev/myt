package repository

import "github.com/semirm-dev/myt/product"

func inMemoryToProduct(p *productInMemory) *product.Product {
	return &product.Product{
		Sku:      p.Sku,
		Name:     p.Name,
		Category: p.Category,
		Price: &product.Price{
			Original: p.Price,
		},
	}
}

func inMemoryToProducts(inMemProducts []*productInMemory) []*product.Product {
	var productsResp []*product.Product

	for _, p := range inMemProducts {
		productsResp = append(productsResp, inMemoryToProduct(p))
	}

	return productsResp
}
