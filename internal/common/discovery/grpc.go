package discovery

import (
	"context"
	"fmt"
	"math/rand"

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

func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("got empty %s addrs from consul", serviceName)
	}
	i := rand.Intn(len(addrs))
	logrus.Infof("Discovered %d instance of %s, addrs=%v", len(addrs), serviceName, addrs)
	return addrs[i], nil
}
