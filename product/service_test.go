package product_test

import (
	"context"
	"github.com/semirm-dev/myt/product"
	pbProduct "github.com/semirm-dev/myt/product/proto"
	"github.com/semirm-dev/myt/product/repository"
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
	addr    = "8001"
)

var lis *bufconn.Listener

type mockedDiscountProvider struct {
	Addr string
}

func (provider *mockedDiscountProvider) ApplyDiscount(ctx context.Context, products []*pbProduct.ProductMessage) ([]*pbProduct.ProductMessage, error) {
	return products, nil
}

func init() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	pbProduct.RegisterProductServer(srv, product.NewService(addr, repository.NewInMemoryRepository(), &mockedDiscountProvider{}))

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

func grpcClient() pbProduct.ProductClient {
	conn := grpcConn(addr)
	return pbProduct.NewProductClient(conn)
}

func TestDefaultService_GetProductsByFilter(t *testing.T) {
	client := grpcClient()

	testSuite := map[string]struct {
		category      string
		plt           int
		expectedCount int
	}{
		"all products": {
			category:      "",
			plt:           0,
			expectedCount: 5,
		},
		"boots category": {
			category:      "boots",
			plt:           0,
			expectedCount: 3,
		},
	}

	for name, ts := range testSuite {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetProductsByFilter(context.Background(), &pbProduct.GetProductsByFilterRequest{
				ByCategory:    ts.category,
				PriceLessThen: int64(ts.plt),
			})

			assert.NoError(t, err)
			assert.Equal(t, ts.expectedCount, len(resp.Products))
		})
	}
}
