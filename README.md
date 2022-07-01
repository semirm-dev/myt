**Run tests**
```shell
go test ./... -v
```

**Run all services**
```shell
docker-compose up
```

**Product service**
* Runs on port 8001 (can be changed)
* Main service to get products
* Communicates to Discounts service through grpc to apply discounts on given products

**Discount service**
* Runs on port 8002 (can be changed)
* Mainly used to apply discounts on given products

**Gateway**
* Runs on port 8000 (can be changed)
* Exposes HTTP API for services (Product only in this case)
* GET localhost:8000/products?priceLessThan=89000&category=boots
* Basic auth needed (can be changed):
    * Username: default
    * Password: default