package config

import (
	"os"
	"sync"

	"github.com/spf13/viper"
)

// 配置结构体
type Config struct {
	ApiBase  ApiBase  `mapstructure:"api-base"`
	Database Database `mapstructure:"database"`
	JWT      JWT      `mapstructure:"jwt"`
}

type ApiBase struct {
	ListenPort int `mapstructure:"listenPort"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWT struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

var (
	Conf *Config
	once sync.Once
)

func Init() {
	once.Do(func() {
		viper.SetConfigType("toml")

		// 根据环境选择配置文件
		env := os.Getenv("GO_ENV")
		if env == "test" {
			viper.SetConfigName("test")
		} else {
			viper.SetConfigName("api")
		}

		viper.AddConfigPath("../config") // 添加上级目录的配置路径
		viper.AddConfigPath("./config")  // 添加当前目录的配置路径

		if err := viper.ReadInConfig(); err != nil {
			// 在测试环境中，如果找不到配置文件，使用默认配置
			if env == "test" {
				Conf = getDefaultConfig()
				return
			}
			panic(err)
		}

		Conf = new(Config)
		if err := viper.Unmarshal(Conf); err != nil {
			panic(err)
		}
	})
}

func getDefaultConfig() *Config {
	return &Config{
		JWT: JWT{
			Secret:     "test-secret-key",
			Expiration: 3600,
		},
		Database: Database{
			Host:     "127.0.0.1",
			Port:     3307,
			User:     "root",
			Password: "200455",
			DBName:   "blog",
		},
	}
}

func GetConfig() *Config {
	if Conf == nil {
		Init()
	}
	return Conf
}
