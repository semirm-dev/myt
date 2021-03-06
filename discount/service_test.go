package discount_test

import (
	"context"
	"github.com/semirm-dev/myt/discount"
	pbDiscount "github.com/semirm-dev/myt/discount/proto"
	"github.com/semirm-dev/myt/discount/repository"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

const (
	bufSize = 1024 * 1024
	addr    = "8002"
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	pbDiscount.RegisterDiscountProviderServer(srv, discount.NewService(addr, repository.NewInMemoryRepository()))

	go func() {
		if err := srv.Serve(lis); err != nil {
			logrus.Fatalf("grpc server failed: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func grpcConn(addr string) *grpc.ClientConn {
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bufDialer),
	}

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		logrus.Fatal(err)
	}

	return conn
}

func grpcClient() pbDiscount.DiscountProviderClient {
	conn := grpcConn(addr)
	return pbDiscount.NewDiscountProviderClient(conn)
}

func TestDefaultService_ApplyDiscount(t *testing.T) {
	testSuite := map[string]struct {
		product                    *pbProduct.ProductMessage
		expectedPriceOriginal      int64
		expectedPriceFinal         int64
		expectedDiscountPercentage string
	}{
		"given sku with highest discount of 15% and lowest discount of 5% will apply 15%": {
			product: &pbProduct.ProductMessage{
				Sku:      "000003",
				Name:     "my product",
				Category: "my category",
				Price: &pbProduct.PriceMessage{
					Original: 1000,
				},
			},
			expectedPriceOriginal:      1000,
			expectedPriceFinal:         850,
			expectedDiscountPercentage: "15%",
		},
		"given category with highest discount of 30% and lowest discount of 15% will apply 30%": {
			product: &pbProduct.ProductMessage{
				Sku:      "000001",
				Name:     "my product",
				Category: "boots",
				Price: &pbProduct.PriceMessage{
					Original: 1000,
				},
			},
			expectedPriceOriginal:      1000,
			expectedPriceFinal:         700,
			expectedDiscountPercentage: "30%",
		},
		"given product without applicable discount": {
			product: &pbProduct.ProductMessage{
				Sku:      "000004",
				Name:     "my product",
				Category: "sandals",
				Price: &pbProduct.PriceMessage{
					Original: 1000,
				},
			},
			expectedPriceOriginal:      1000,
			expectedPriceFinal:         1000,
			expectedDiscountPercentage: "",
		},
	}

	client := grpcClient()

	for name, ts := range testSuite {
		t.Run(name, func(t *testing.T) {
			resp, err := client.ApplyDiscount(context.Background(), &pbDiscount.DiscountsRequest{
				Products: []*pbProduct.ProductMessage{
					ts.product,
				},
			})

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, 1, len(resp.Products))

			product := resp.Products[0]
			assert.NotNil(t, product.Price)
			assert.Equal(t, ts.expectedPriceOriginal, product.Price.Original)
			assert.Equal(t, ts.expectedPriceFinal, product.Price.Final)
			assert.Equal(t, ts.expectedDiscountPercentage, product.Price.DiscountPercentage)
		})
	}
}
