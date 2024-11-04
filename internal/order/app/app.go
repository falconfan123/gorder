package app

import "github.com/falconfan123/gorder/order/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
