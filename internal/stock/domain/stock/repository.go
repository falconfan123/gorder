package stock

import (
	"context"
	"fmt"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string)
}

type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("these items not found in stock:%s", strings.Join(e.Missing, ","))
}
