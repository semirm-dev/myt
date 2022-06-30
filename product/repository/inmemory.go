package repository

import (
	"context"
	"github.com/semirm-dev/myt/product"
	"strings"
)

var productsInMemory = []*productInMemory{
	{
		Sku:      "000001",
		Name:     "BV Lean leather ankle boots",
		Category: "boots",
		Price:    89000,
	},
	{
		Sku:      "000002",
		Name:     "BV Lean leather ankle boots",
		Category: "boots",
		Price:    99000,
	},
	{
		Sku:      "000003",
		Name:     "Ashlington leather ankle boots",
		Category: "boots",
		Price:    71000,
	},
	{
		Sku:      "000004",
		Name:     "Naima embellished suede sandals",
		Category: "sandals",
		Price:    79500,
	},
	{
		Sku:      "000005",
		Name:     "Nathane leather sneakers",
		Category: "sneakers",
		Price:    59000,
	},
}

type productInMemory struct {
	Sku      string
	Name     string
	Category string
	Price    int
}

// inMemoryRepository is in memory data store implementation for products
type inMemoryRepository struct {
	products []*productInMemory
}

func NewInMemoryRepository() product.Repository {
	return &inMemoryRepository{
		products: productsInMemory,
	}
}

func (repo *inMemoryRepository) GetProductsByFilter(ctx context.Context, filter *product.Filter) ([]*product.Product, error) {
	products := repo.products

	if strings.TrimSpace(filter.ByCategory) != "" {
		products = filterByCategory(products, filter.ByCategory)
	}

	if filter.PriceLessThan > 0 {
		products = filterByPriceLessThan(products, filter.PriceLessThan)
	}

	return inMemoryToProducts(products), nil
}

// filterByCategory will get products with a given category.
// Where condition in database query would solve this issue.
func filterByCategory(products []*productInMemory, category string) []*productInMemory {
	var filtered []*productInMemory

	for _, p := range products {
		if p.Category == category {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

// filterByPriceLessThan will get products with price less than given price.
// Where condition in database query would solve this issue.
func filterByPriceLessThan(products []*productInMemory, price int) []*productInMemory {
	var filtered []*productInMemory

	for _, p := range products {
		if p.Price <= price {
			filtered = append(filtered, p)
		}
	}

	return filtered
}
