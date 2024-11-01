package main

import (
	"github.com/falconfan123/gorder/common/genproto/stockpb"
	"github.com/falconfan123/gorder/common/server"
	"github.com/falconfan123/gorder/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service_name")
	serverType := viper.Get("stock.server-to-run")

	switch serverType {
	case "grpc":
		server.RunGRPCServe(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		//"TODO"
	default:
		panic("unexpected server type")

	}
}
