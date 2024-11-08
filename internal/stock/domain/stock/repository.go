package stock

import (
	"context"
	"fmt"
	"github.com/falconfan123/gorder/common/genproto/orderpb"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("these items not found in stock:%s", strings.Join(e.Missing, ","))
}
