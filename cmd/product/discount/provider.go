package discount

import (
	"context"
	pbDiscount "github.com/semirm-dev/myt/discount/proto"
	"github.com/semirm-dev/myt/internal/grpc"
)

type Provider struct {
	Addr string
}

func (provider *Provider) GetDiscounts(ctx context.Context) (*pbDiscount.DiscountsResponse, error) {
	conn := grpc.CreateClientConnection(provider.Addr)
	discountsClient := pbDiscount.NewDiscountProviderClient(conn)

	return discountsClient.GetDiscounts(ctx, &pbDiscount.DiscountsRequest{})
}
