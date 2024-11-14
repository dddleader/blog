package article

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

type Article struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Summary string   `yaml:"summary"`
	Cover   string   `yaml:"cover"`
	Tags    []string `yaml:"tags"`
	Content string   `yaml:"-"`
	Path    string   `yaml:"-"`
}

type ArticleList struct {
	Articles []Article
	Total    int
	Page     int
	PageSize int
}

func LoadArticles(page, pageSize int) (*ArticleList, error) {
	files, err := ioutil.ReadDir("articles")
	if err != nil {
		return nil, err
	}

	var articles []Article
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			article, err := parseArticle(file.Name())
			if err != nil {
				continue
			}
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
	content, err := ioutil.ReadFile(filepath.Join("articles", filename))
	if err != nil {
		return Article{}, err
	}

	// 分割前置元数据和内容
	parts := strings.Split(string(content), "---\n")
	if len(parts) < 3 {
		return Article{}, fmt.Errorf("invalid article format")
	}

	var article Article
	err = yaml.Unmarshal([]byte(parts[1]), &article)
	if err != nil {
		return Article{}, err
	}

	// 解析Markdown内容
	m := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	var buf strings.Builder
	if err := m.Convert([]byte(parts[2]), &buf); err != nil {
		return Article{}, err
	}
	article.Content = buf.String()
	article.Path = strings.TrimSuffix(filename, ".md")

	return article, nil
}

func GetArticle(path string) (*Article, error) {
	article, err := parseArticle(path + ".md")
	if err != nil {
		return nil, err
	}
	return &article, nil
}
