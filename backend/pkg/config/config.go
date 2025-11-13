package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()

	// 设置环境变量前缀为 RS（Robot Shop）
	conf.SetEnvPrefix("RS")

	// 启用环境变量自动读取
	conf.AutomaticEnv()

	// 尝试读取配置文件（如果文件不存在，仍可从环境变量读取）
	if path != "" {
		conf.SetConfigFile(path)
		err := conf.ReadInConfig()
		if err != nil {
			fmt.Printf("Warning: failed to read config file %s: %v\n", path, err)
			fmt.Println("Will use environment variables with RS_ prefix")
		} else {
			fmt.Printf("Loaded config file: %s\n", path)
		}
	}

	return conf
}
