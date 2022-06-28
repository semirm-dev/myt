package main

import (
	"flag"
	"github.com/semirm-dev/myt/gateway"
	"github.com/semirm-dev/myt/gateway/handlers"
	"github.com/semirm-dev/myt/internal/web"
	"github.com/semirm-dev/myt/product"
)

var (
	httpAddr    = flag.String("http", ":8000", "Http address")
	authUser    = flag.String("usr", "default", "default username for basic auth")
	authPwd     = flag.String("pwd", "default", "default password for basic auth")
	productsUri = flag.String("products_uri", ":8001", "Products Service address")
)

func main() {
	flag.Parse()

	router := web.NewRouter()

	router.Use(gateway.BasicAuth(*authUser, *authPwd))

	productsClient := product.NewClient(*productsUri)
	router.GET("products", handlers.GetProducts(productsClient))

	web.ServeHttp(*httpAddr, "gateway", router)
}
