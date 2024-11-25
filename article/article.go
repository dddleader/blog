package article

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"blog/models"
)

// @Summary 获取文章列表
// @Description 获取所有已发布的文章列表，支持分页
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} models.ArticleListResponse
// @Router /articles [get]
func GetArticles(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 修改数据库连接字符串
	const dbConnString = "root:200455@tcp(127.0.0.1:3307)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dbConnString)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	// 获取��数
	var total int64
	err = db.QueryRow("SELECT COUNT(*) FROM articles WHERE on_show = true").Scan(&total)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章总数失败"})
		return
	}

	// 查询文章列表
	offset := (page - 1) * pageSize
	rows, err := db.Query(`
		SELECT id, title, content, summary, cover, created_at, updated_at, status, views, tags, on_show
		FROM articles 
		WHERE on_show = true 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章列表失败"})
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var tagsJSON []byte

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Summary,
			&article.Cover,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Status,
			&article.Views,
			&tagsJSON,
			&article.OnShow,
		)
		if err != nil {
			logrus.Error("扫描文章数据失败:", err)
			continue
		}

		// 解析tags JSON
		if err := json.Unmarshal(tagsJSON, &article.Tags); err != nil {
			logrus.Error("解析标签失败:", err)
			article.Tags = []string{}
		}

		articles = append(articles, article)
	}

	// 检查遍历错误
	if err = rows.Err(); err != nil {
		logrus.Error("遍历文章数据失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文章数据失败"})
		return
	}

	// 添加日志
	logrus.Infof("查询到 %d 篇文章", len(articles))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": models.ArticleListResponse{
			Articles: articles,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// @Summary 获取单篇文章
// @Description 根据ID获取单篇文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} models.Article
// @Router /articles/{id} [get]
func GetSingleArticle(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("mysql", "root:200455@tcp(127.0.0.1:3307)/blog")
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	var article models.Article
	var tagsJSON []byte
	err = db.QueryRow(`
		SELECT id, title, content, summary, cover, created_at, updated_at, status, views, tags, on_show
		FROM articles 
		WHERE id = ? AND on_show = true`,
		id,
	).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Summary,
		&article.Cover,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Status,
		&article.Views,
		&tagsJSON,
		&article.OnShow,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
		return
	}

	// 解析tags JSON
	if err := json.Unmarshal(tagsJSON, &article.Tags); err != nil {
		logrus.Error(err)
		article.Tags = []string{}
	}

	// 更新浏览次数
	_, err = db.Exec("UPDATE articles SET views = views + 1 WHERE id = ?", id)
	if err != nil {
		logrus.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": article,
	})
}
