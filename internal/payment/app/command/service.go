package command

import (
	"context"

	"github.com/falconfan123/gorder/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
