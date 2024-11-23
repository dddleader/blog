package article

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminLogin(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 测试用例
	tests := []struct {
		name       string
		reqBody    LoginRequest
		wantStatus int
		wantError  bool
	}{
		{
			name: "成功登录",
			reqBody: LoginRequest{
				Username: "admin",
				Password: "200455",
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "用户名不存在",
			reqBody: LoginRequest{
				Username: "nonexistent",
				Password: "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  true,
		},
		{
			name: "密码错误",
			reqBody: LoginRequest{
				Username: "admin",
				Password: "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
			wantError:  true,
		},
		{
			name: "请求参数为空",
			reqBody: LoginRequest{
				Username: "",
				Password: "",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.POST("/admin/login", AdminLogin)

			// 创建请求体
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/admin/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 记录响应
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 断言状态码
			assert.Equal(t, tt.wantStatus, w.Code)

			// 解析响应
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 验证响应内容
			if tt.wantError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			}
		})
	}
}

// 测试数据库相关函数
func TestGetUserByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "存在的用户",
			username: "admin",
			wantErr:  false,
		},
		{
			name:     "不存在的用户",
			username: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := getUserByUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
			}
		})
	}
}

// Mock数据库测试
func TestAdminLoginWithMockDB(t *testing.T) {
	// TODO: 使用 sqlmock 模拟数据库
	// 这需要重构 getUserByUsername 函数以接受数据库连接作为参数
}
