package jwt

import (
	"blog/config"
)

// 初始化测试配置
func setupTestConfig() {
	config.Conf = &config.Config{
		JWT: config.JWT{
			Secret:     "test-secret-key",
			Expiration: 3600,
		},
	}
}
