package product

import "context"

type Repository interface {
	GetProductsByFilter(ctx context.Context, filter *Filter) ([]*Product, error)
}
