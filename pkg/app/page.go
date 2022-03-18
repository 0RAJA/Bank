package app

import (
	"github.com/gin-gonic/gin"
)

//分页处理

var (
	DefaultPageSize int
	MaxPageSize     int
	PageKey         string //url中page关键字
	PageSizeKey     string //pagesize关键字
)

// Init 初始化默认页数大小和最大页数限制以及查询的关键字
func Init(defaultPageSize, maxPageSize int, pageKey, pageSizeKey string) {
	DefaultPageSize = defaultPageSize
	MaxPageSize = maxPageSize
	PageKey = pageKey
	PageSizeKey = pageSizeKey
}

type Pager struct {
	Page      int `json:"page,omitempty"`
	PageSize  int `json:"page_size,omitempty"`
	TotalRows int `json:"total_rows"`
}

func NewPager(page int, pageSize int) *Pager {
	return &Pager{Page: page, PageSize: pageSize}
}

// GetPage 获取页数
func GetPage(c *gin.Context) int {
	page := StrTo(c.Query(PageKey)).MustInt()
	if page <= 0 {
		return 1
	}
	return page
}

// GetPageSize 获取pageSize
func GetPageSize(c *gin.Context) int {
	pageSize := StrTo(c.Query(PageSizeKey)).MustInt()
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

// GetPageOffset 通过page和pageSize获取偏移值
func GetPageOffset(page, pageSize int) (result int) {
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return
}
