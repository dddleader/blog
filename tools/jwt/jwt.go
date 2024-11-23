package jwt

import (
	"errors"
	"time"

	"blog/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims结构
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT服务
type JWTService struct {
	secretKey string
	duration  time.Duration
}

// 创建JWT服务
func NewJWTService() *JWTService {
	// 如果配置为空，使用默认值
	secret := "default-secret-key"
	duration := 24 * 60 * 60 // 24小时

	if config.Conf != nil && config.Conf.JWT.Secret != "" {
		secret = config.Conf.JWT.Secret
		duration = config.Conf.JWT.Expiration
	}

	return &JWTService{
		secretKey: secret,
		duration:  time.Duration(duration) * time.Second,
	}
}

// 生成token
func (j *JWTService) GenerateToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-system",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// 验证token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
