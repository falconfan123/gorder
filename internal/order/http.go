package main

import (
	"fmt"
	"github.com/falconfan123/gorder/common"
	client "github.com/falconfan123/gorder/common/client/order"
	"github.com/falconfan123/gorder/order/app"
	"github.com/falconfan123/gorder/order/app/command"
	"github.com/falconfan123/gorder/order/app/dto"
	"github.com/falconfan123/gorder/order/app/query"
	"github.com/falconfan123/gorder/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	//TODO implement me
	panic("implement me")
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {
	//ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
	//defer span.End()

	var (
		req  client.CreateOrderRequest
		err  error
		resp dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResponse{
		CustomerID:  req.CustomerId,
		OrderID:     r.OrderID,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp interface{}
	)
	defer func() {
		H.Response(c, err, resp)
	}()
	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	//between app and http =====> dto
	if err != nil {
		return
	}
	resp = convertor.NewOrderConvertor().EntityToClient(o)
}
