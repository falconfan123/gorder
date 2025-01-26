package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() error {
	relPath, err := GetRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relPath)
	viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	return viper.ReadInConfig()
}

func GetRelativePathFromCaller() (relPath string, err error) {
	callerpwd, err := os.Getwd()
	if err != nil {
		return
	}
	//获取当前这个文件的目录
	_, here, _, _ := runtime.Caller(0)
	relPath, err = filepath.Rel(callerpwd, filepath.Dir(here)) //从哪到哪的相对路径
	fmt.Printf("caller from : %s, here : %s, relPath : %s", callerpwd, here, relPath)
	return
}
