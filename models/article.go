package models

import "time"

// 文章基础结构
type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Cover     string    `json:"cover"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `json:"status"`
	Views     uint      `json:"views"`
	Tags      []string  `json:"tags"`
	OnShow    bool      `json:"on_show"`
}

// 文章列表请求参数
type ArticleListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=50"`
	Tag      string `form:"tag"`
	Search   string `form:"search"`
}

// 文��列表响应
type ArticleListResponse struct {
	Articles []Article `json:"articles"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}
