package discovery

import (
	"context"
	"github.com/falconfan123/gorder/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 返回cleanup函数和error
func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.address"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	//获得grpc地址
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	logrus.WithFields(logrus.Fields{
		"grpcAddr": grpcAddr,
	}).Info("in RegisterToConsul")
	//把grpc地址注册到consul内
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}
	//启动协程去监视心跳
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
		}
	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("register to consul")
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}
