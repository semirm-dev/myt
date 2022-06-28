package discount

import pbDiscount "github.com/semirm-dev/myt/discount/proto"

func toProtoDiscountMessage(discount *Discount) *pbDiscount.DiscountMessage {
	return &pbDiscount.DiscountMessage{
		Sku:                discount.Sku,
		Category:           discount.Category,
		DiscountPercentage: int64(discount.Percentage),
	}
}

func toProtoDiscountsResponse(discounts []*Discount) *pbDiscount.DiscountsResponse {
	var resp []*pbDiscount.DiscountMessage

	for _, discount := range discounts {
		resp = append(resp, toProtoDiscountMessage(discount))
	}

	return &pbDiscount.DiscountsResponse{
		Discounts: resp,
	}
}
