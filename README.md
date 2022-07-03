**Run tests**
```shell
go test ./... -v
```

**Run all services**
```shell
docker-compose up
```

**Product service**
* Runs on port 8001 (configurable)
* Main service to get products
* Communicates to Discounts service through grpc to apply discounts on given products

**Discount service**
* Runs on port 8002 (configurable)
* Used to apply discounts on given products

**Gateway**
* Runs on port 8000 (configurable)
* Exposes HTTP API for services (Product only in this case)
* Endpoint: _GET localhost:8000/products?priceLessThan=89000&category=boots_
* Basic auth is needed (configurable):
    * Username: default
    * Password: default

> All repositories (data stores) are implemented as in memory for the sake of the simplicity.

**Tasks**
- [x] Implement services: _Product, Discount_
- [x] Implement _Gateway_ to expose services (Product only in this case)
- [x] Write tests (Product and Discount services)
- [ ] Write tests for API Gateway
- [x] All services and tests runnable by 1 command
- [x] Products in _boots_ category have discount of 30%
- [x] Products with sku _000003_ have discount of 15%
- [x] Always apply the biggest discount
- [x] Filter products by _category_
- [x] Filter products by _price less than_
- [x] If product doesnt have discount then price.Final == price.Original
- [x] Default currency is _EUR_
- [ ] Implement pagination (by default 5 items per page)
- [ ] Provide SQL implementation instead of in memory 