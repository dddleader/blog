package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog/config"
	"blog/models"
	"blog/tools/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化测试配置
	config.Conf = &config.Config{
		JWT: config.JWT{
			Secret:     "test-secret-key",
			Expiration: 3600,
		},
		Database: config.Database{
			Host:     "127.0.0.1",
			Port:     3307,
			User:     "root",
			Password: "200455",
			DBName:   "blog",
		},
	}
}

// 辅助函数：获取测试用JWT token
func getTestToken(t *testing.T) string {
	jwtService := jwt.NewJWTService()
	token, err := jwtService.GenerateToken(1, "admin")
	assert.NoError(t, err)
	return token
}

func TestCreateArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		article    models.Article
		withAuth   bool
		wantStatus int
	}{
		{
			name: "成功创建文章",
			article: models.Article{
				Title:   "测试文章",
				Content: "测试内容",
				Summary: "测试摘要",
				Status:  1,
				Tags:    []string{"测试", "Go"},
				OnShow:  true,
			},
			withAuth:   true,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/articles", CreateArticle)

			body, _ := json.Marshal(tt.article)
			req := httptest.NewRequest("POST", "/articles", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.withAuth {
				req.Header.Set("Authorization", "Bearer "+getTestToken(t))
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "id")
				assert.Contains(t, response, "message")
			}
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		articleID  string
		article    models.Article
		withAuth   bool
		wantStatus int
	}{
		{
			name:      "成功更新文章",
			articleID: "1",
			article: models.Article{
				Title:   "更新的标题",
				Content: "更新的内容",
			},
			withAuth:   true,
			wantStatus: http.StatusOK,
		},
		{
			name:      "文章不存在",
			articleID: "999",
			article: models.Article{
				Title: "不存在",
			},
			withAuth:   true,
			wantStatus: http.StatusNotFound,
		},
		{
			name:      "未授权更新",
			articleID: "1",
			article: models.Article{
				Title: "未授权",
			},
			withAuth:   false,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.PUT("/articles/:id", UpdateArticle)

			url := fmt.Sprintf("/articles/%s", tt.articleID)
			body, _ := json.Marshal(tt.article)
			req := httptest.NewRequest("PUT", url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.withAuth {
				req.Header.Set("Authorization", "Bearer "+getTestToken(t))
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		articleID  string
		withAuth   bool
		wantStatus int
	}{
		{
			name:       "成功删除文章",
			articleID:  "1",
			withAuth:   true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "文章不存在",
			articleID:  "999",
			withAuth:   true,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "未授权删除",
			articleID:  "1",
			withAuth:   false,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.DELETE("/articles/:id", DeleteArticle)

			url := fmt.Sprintf("/articles/%s", tt.articleID)
			req := httptest.NewRequest("DELETE", url, nil)

			if tt.withAuth {
				req.Header.Set("Authorization", "Bearer "+getTestToken(t))
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// 使用 sqlmock 测试数据库操作
func TestArticleWithMockDB(t *testing.T) {
	// TODO: 添加使用 sqlmock 的数据库操作测试
	// 这需要重构数据库操作，使其接受 db 参数而不是在函数内部创建连接
}
