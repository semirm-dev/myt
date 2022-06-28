package product

type Product struct {
	Sku      string
	Name     string
	Category string
	Price    *Price
}

type Price struct {
	Original           int
	Final              int
	DiscountPercentage string
	Currency           string
}
