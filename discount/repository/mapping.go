package repository

import "github.com/semirm-dev/myt/discount"

func inMemoryToDiscount(d *discountInMemory) *discount.Discount {
	return &discount.Discount{
		Sku:        d.sku,
		Category:   d.category,
		Percentage: d.amount,
	}
}

func inMemoryToDiscounts(inMemDiscounts []*discountInMemory) []*discount.Discount {
	var resp []*discount.Discount

	for _, d := range inMemDiscounts {
		resp = append(resp, inMemoryToDiscount(d))
	}

	return resp
}
