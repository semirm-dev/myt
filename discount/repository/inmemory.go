package repository

import (
	"context"
	"github.com/semirm-dev/myt/discount"
)

type discountInMemory struct {
	sku      string
	category string
	amount   int
}

var discountsInMemory = []*discountInMemory{
	{
		category: "boots",
		amount:   30,
	},
	{
		sku:    "000003",
		amount: 15,
	},
	{
		category: "boots",
		amount:   50,
	},
}

type inMemoryRepository struct {
	discounts []*discountInMemory
}

func NewInMemoryRepository() discount.Repository {
	return &inMemoryRepository{
		discounts: discountsInMemory,
	}
}

func (repo *inMemoryRepository) GetDiscount(ctx context.Context) ([]*discount.Discount, error) {
	discounts := repo.discounts

	return inMemoryToDiscounts(discounts), nil
}
