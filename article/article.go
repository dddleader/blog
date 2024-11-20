package article

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

type Article struct {
	ID      string   `yaml:"id"` // 新增ID字段
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Summary string   `yaml:"summary"`
	Cover   string   `yaml:"cover"`
	Tags    []string `yaml:"tags"`
	Content string   `yaml:"-"`
	Path    string   `yaml:"-"`
}
type ArticleList struct {
	Articles []Article `json:"articles"` // 添加 json tag
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

// 获取文章列表
func GetArticles(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 获取文章列表
	articles, err := LoadArticles(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to load articles",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"articles": articles.Articles,
			"total":    articles.Total,
			"page":     articles.Page,
			"pageSize": articles.PageSize,
		},
	})
}
func LoadArticles(page, pageSize int) (*ArticleList, error) {
	files, err := os.ReadDir("articles")
	if err != nil {
		return nil, err
	}

	var articles []Article
	for _, file := range files {
		// 只读取md文件
		if filepath.Ext(file.Name()) == ".md" {
			article, err := parseArticle(file.Name())
			if err != nil {
				logrus.Errorf("Error parsing article %s: %v", file.Name(), err)
				continue
			}
			logrus.Infof("Parsed article: %s", article.Title)
			articles = append(articles, article)
		}
	}

	// 计算分页
	total := len(articles)
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	if start >= total {
		start = total
	}

	return &ArticleList{
		Articles: articles[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func parseArticle(filename string) (Article, error) {
	content, err := os.ReadFile(filepath.Join("articles", filename))
	if err != nil {
		return Article{}, err
	}

	// 修改这里：使用更精确的分割方式
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return Article{}, fmt.Errorf("invalid article format")
	}

	// 去除前后空白
	yamlContent := strings.TrimSpace(parts[1])
	markdownContent := strings.TrimSpace(parts[2])

	var article Article
	err = yaml.Unmarshal([]byte(yamlContent), &article)
	if err != nil {
		return Article{}, fmt.Errorf("yaml parse error: %v", err)
	}

	// 解析Markdown内容
	m := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	var buf strings.Builder
	if err := m.Convert([]byte(markdownContent), &buf); err != nil {
		return Article{}, err
	}
	article.Content = buf.String()
	article.Path = strings.TrimSuffix(filename, ".md")

	return article, nil
}

func GetSingleArticle(c *gin.Context) {
	// 从URL参数中获取path
	path := c.Param("path")

	// 获取文章
	article, err := parseArticle(path + ".md")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Article not found",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": article,
	})
}
