package discount

import (
	"context"
	pbDiscount "github.com/semirm-dev/myt/discount/proto"
	"github.com/semirm-dev/myt/internal/grpc"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "discount service"

func NewService(addr string, repo Repository) *defaultService {
	return &defaultService{
		addr: addr,
		repo: repo,
	}
}

// defaultService will expose discount service via grpc
type defaultService struct {
	pbDiscount.UnimplementedDiscountProviderServer
	addr string
	repo Repository
}

func (svc *defaultService) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.addr, serviceName)
}

func (svc *defaultService) RegisterGrpcServer(server *grpcLib.Server) {
	pbDiscount.RegisterDiscountProviderServer(server, svc)
}

func (svc *defaultService) GetDiscounts(ctx context.Context, req *pbDiscount.DiscountsRequest) (*pbDiscount.DiscountsResponse, error) {
	discounts, err := svc.repo.GetDiscount(ctx)
	if err != nil {
		return nil, err
	}

	return toProtoDiscountsResponse(discounts), nil
}
