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
		sku:    "000003",
		amount: 35,
	},
	{
		sku:    "000003",
		amount: 25,
	},
	{
		category: "boots",
		amount:   50,
	},
	{
		category: "boots",
		amount:   40,
	},
	{
		sku:    "000010",
		amount: 25,
	},
	{
		sku:    "000015",
		amount: 25,
	},
}

// inMemoryRepository is in memory data store implementation for discounts
type inMemoryRepository struct {
	discounts []*discountInMemory
}

func NewInMemoryRepository() discount.Repository {
	return &inMemoryRepository{
		discounts: discountsInMemory,
	}
}

func (repo *inMemoryRepository) GetDiscounts(ctx context.Context, filter *discount.Filter) ([]*discount.Discount, error) {
	var discounts []*discountInMemory

	// solved with Where conditions in database query

	for _, s := range filter.Sku {
		for _, d := range repo.discounts {
			if d.sku == s {
				discounts = append(discounts, d)
			}
		}
	}

	for _, c := range filter.Category {
		for _, d := range repo.discounts {
			if d.category == c {
				discounts = append(discounts, d)
			}
		}
	}

	return inMemoryToDiscounts(discounts), nil
}
