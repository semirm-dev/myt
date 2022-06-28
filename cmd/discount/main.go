package main

import (
	"context"
	"flag"
	"github.com/semirm-dev/myt/discount"
	"github.com/semirm-dev/myt/discount/repository"
)

var (
	addr = flag.String("addr", ":8002", "Discounts Service address")
)

func main() {
	flag.Parse()

	svc := discount.NewService(*addr, repository.NewInMemoryRepository())

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	svc.ListenForConnections(rootCtx)
}
