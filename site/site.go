/**
 * 静态文件服务
 */
package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"blog/article"
	"blog/config"
)

type Site struct {
}

func New() *Site {
	return &Site{}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("./site/index.html")
	_, _ = fmt.Fprintf(w, string(data))
	return
}

func server(fs http.FileSystem) http.Handler {
	// 创建文件服务器
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取文件路径
		filePath := path.Clean("./site" + r.URL.Path)
		// 检查文件是否存在
		_, err := os.Stat(filePath)
		if err != nil {
			notFound(w, r)
			return
		}
		// 返回文件
		fileServer.ServeHTTP(w, r)
	})
}

func (s *Site) Run() {
	siteConfig := config.GetConfig().Site
	port := siteConfig.SiteBase.ListenPort

	mux := http.NewServeMux()

	// 静态文件服务
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("site/static"))))

	// API endpoints
	mux.HandleFunc("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		page := 1
		pageSize := 10

		if p := r.URL.Query().Get("page"); p != "" {
			fmt.Sscanf(p, "%d", &page)
		}
		if ps := r.URL.Query().Get("pageSize"); ps != "" {
			fmt.Sscanf(ps, "%d", &pageSize)
		}

		articles, err := article.LoadArticles(page, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	})

	mux.HandleFunc("/api/article/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/article/")
		article, err := article.GetArticle(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	})

	// 所有其他路由返回index.html
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/index.html")
	})

	addr := fmt.Sprintf(":%d", port)
	logrus.Info("starting site server on", addr)
	logrus.Fatal(http.ListenAndServe(addr, mux))
}
