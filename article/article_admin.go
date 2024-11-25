package article

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"blog/models"
)

// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param article body models.CreateArticleRequest true "文章信息"
// @Success 200 {object} map[string]interface{}
// @Router /admin/articles [post]
func CreateArticle(c *gin.Context) {
	var req models.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	db, err := sql.Open("mysql", "root:200455@tcp(127.0.0.1:3307)/blog?charset=utf8mb4&parseTime=True")
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	// 将tags转换为JSON
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标签格式错误"})
		return
	}

	// 插入文章
	result, err := db.Exec(`
		INSERT INTO articles (
			title, content, summary, cover, status, tags, 
			on_show, views, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, true, 0, NOW(), NOW())`,
		req.Title, req.Content, req.Summary, req.Cover,
		req.Status, tagsJSON,
	)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"message": "文章创建成功",
		"id":      id,
	})
}

// @Summary 更新文章
// @Description 更新现有文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "文章ID"
// @Param article body models.UpdateArticleRequest true "文章信息"
// @Success 200 {object} map[string]interface{}
// @Router /admin/articles/{id} [put]
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	db, err := sql.Open("mysql", "root:200455@tcp(127.0.0.1:3307)/blog?charset=utf8mb4&parseTime=True")
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标签格式错误"})
		return
	}

	// 更新所有字段，包括 cover
	_, err = db.Exec(`
		UPDATE articles 
		SET title=?, content=?, summary=?, cover=?, status=?, tags=?
		WHERE id=?`,
		req.Title, req.Content, req.Summary, req.Cover,
		req.Status, tagsJSON, id,
	)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章更新成功"})
}

// @Summary 删除文章
// @Description 软删除文章（设置为不显示）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("mysql", "root:200455@tcp(127.0.0.1:3307)/blog?charset=utf8mb4&parseTime=True")
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE articles SET on_show=false WHERE id=?", id)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}
