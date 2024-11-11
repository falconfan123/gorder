package domain

import (
	"context"
	"github.com/falconfan123/gorder/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}
