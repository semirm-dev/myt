package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/myt/product"
	"github.com/semirm-dev/myt/product/repository"
)

var (
	addr = flag.String("addr", ":8001", "Product Service address")
)

func main() {
	flag.Parse()

	svc := product.NewService(*addr, repository.NewInMemoryRepository())

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	svc.ListenForConnections(rootCtx)
}
