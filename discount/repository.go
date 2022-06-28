package discount

import (
	"context"
)

type Repository interface {
	GetDiscount(ctx context.Context) ([]*Discount, error)
}
