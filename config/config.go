package config

import (
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var Conf *Config

type Config struct {
	ApiBase ApiBase `mapstructure:"api-base"`
	Site    Site    `mapstructure:"site"`
}

type ApiBase struct {
	ListenPort int `mapstructure:"listenPort"`
}

type Site struct {
	SiteBase SiteBase `mapstructure:"site-base"`
}

type SiteBase struct {
	ListenPort int `mapstructure:"listenPort"`
}

func Init() {
	once.Do(func() {
		viper.SetConfigType("toml")
		viper.AddConfigPath("./config")

		// 读取 api 配置
		viper.SetConfigName("api")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		// 读取 site 配置
		viper.SetConfigName("site")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}

		Conf = new(Config)
		err = viper.Unmarshal(Conf)
		if err != nil {
			panic(err)
		}
	})
}

func GetConfig() *Config {
	if Conf == nil {
		Init()
	}
	return Conf
}
