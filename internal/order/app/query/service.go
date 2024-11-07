package query

import (
	"context"
	"github.com/falconfan123/gorder/common/genproto/orderpb"
)

type StockService interface {
	CheckIFItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
