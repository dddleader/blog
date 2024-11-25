package article

import (
	"database/sql"
	"encoding/json"
	"testing"

	"blog/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestInsertArticle(t *testing.T) {
	// 准备测试数据
	article := models.Article{
		Title:     "Go语言 sort 包详解：从基础到进阶",
		Content:   `# Go语言 sort 包详解：从基础到进阶...`, // 这里是完整的文章内容
		Summary:   "详细介绍了 Go 语言 sort 包的使用方法，包括基本排序功能、自定义排序实现以及二分查找等核心功能，适合想要深入了解 Go 排序机制的开发者。",
		Status:    1,
		Tags:      []string{"Go", "标准库", "算法", "排序"},
		OnShow:    true,
		CreatedAt: "2024-01-01 00:00:00",
		UpdatedAt: "2024-01-01 00:00:00",
	}

	// 连接数据库
	dbConnString := "root:200455@tcp(127.0.0.1:3307)/blog?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dbConnString)
	assert.NoError(t, err)
	defer db.Close()

	// 将tags转换为JSON
	tagsJSON, err := json.Marshal(article.Tags)
	assert.NoError(t, err)

	// 插入文章
	result, err := db.Exec(`
		INSERT INTO articles (title, content, summary, cover, status, tags, on_show, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		article.Title, article.Content, article.Summary, article.Cover,
		article.Status, tagsJSON, article.OnShow, article.CreatedAt, article.UpdatedAt,
	)
	assert.NoError(t, err)

	// 验证插入结果
	id, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// 验证文章是否正确插入
	var savedArticle models.Article
	var savedTagsJSON []byte
	err = db.QueryRow(`
		SELECT id, title, content, summary, cover, status, tags, on_show, created_at, updated_at
		FROM articles WHERE id = ?`,
		id,
	).Scan(
		&savedArticle.ID,
		&savedArticle.Title,
		&savedArticle.Content,
		&savedArticle.Summary,
		&savedArticle.Cover,
		&savedArticle.Status,
		&savedTagsJSON,
		&savedArticle.OnShow,
		&savedArticle.CreatedAt,
		&savedArticle.UpdatedAt,
	)
	assert.NoError(t, err)

	// 解析保存的tags
	err = json.Unmarshal(savedTagsJSON, &savedArticle.Tags)
	assert.NoError(t, err)

	// 验证字段值
	assert.Equal(t, article.Title, savedArticle.Title)
	assert.Equal(t, article.Summary, savedArticle.Summary)
	assert.Equal(t, article.Cover, savedArticle.Cover)
	assert.Equal(t, article.Status, savedArticle.Status)
	assert.Equal(t, article.Tags, savedArticle.Tags)
	assert.Equal(t, article.OnShow, savedArticle.OnShow)
	assert.Equal(t, article.CreatedAt, savedArticle.CreatedAt)
	assert.Equal(t, article.UpdatedAt, savedArticle.UpdatedAt)
}
