package main

import (
	"github.com/falconfan123/gorder/common/genproto/orderpb"
	"github.com/falconfan123/gorder/order/app"
	"github.com/falconfan123/gorder/order/app/command"
	"github.com/falconfan123/gorder/order/app/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"customer_iD": req.CustomerID,
		"order_iD":    r.OrderID,
	})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		CustomerID: "fake-customer-id",
		OrderID:    "fake-ID",
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
}
