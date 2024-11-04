package service

import (
	"context"
	"github.com/falconfan123/gorder/common/metrics"
	"github.com/falconfan123/gorder/order/adapters"
	"github.com/falconfan123/gorder/order/app"
	"github.com/falconfan123/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
