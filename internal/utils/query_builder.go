package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QueryParams struct {
	Page         string
	PageSize     string
	CreatedAtGte string
	CreatedAtLte string
	CreatedAtGt  string
	CreatedAtLt  string
	UpdatedAtGte string
	UpdatedAtLte string
	UpdatedAtGt  string
	UpdatedAtLt  string
}

func BuildQuery(c *gin.Context, db *gorm.DB) *gorm.DB {
	params := QueryParams{
		CreatedAtGte: c.Query("created_at__gte"),
		CreatedAtLte: c.Query("created_at__lte"),
		CreatedAtGt:  c.Query("created_at__gt"),
		CreatedAtLt:  c.Query("created_at__lt"),
		UpdatedAtGte: c.Query("updated_at__gte"),
		UpdatedAtLte: c.Query("updated_at__lte"),
		UpdatedAtGt:  c.Query("updated_at__gt"),
		UpdatedAtLt:  c.Query("updated_at__lt"),
		Page:         c.Query("page"),
		PageSize:     c.Query("page_size"),
	}

	query := db

	if params.CreatedAtGte != "" {
		query = query.Where("created_at >= ?", params.CreatedAtGte)
	}
	if params.CreatedAtLte != "" {
		query = query.Where("created_at <= ?", params.CreatedAtLte)
	}
	if params.CreatedAtGt != "" {
		query = query.Where("created_at > ?", params.CreatedAtGt)
	}
	if params.CreatedAtLt != "" {
		query = query.Where("created_at < ?", params.CreatedAtLt)
	}

	if params.UpdatedAtGte != "" {
		query = query.Where("updated_at >= ?", params.UpdatedAtGte)
	}
	if params.UpdatedAtLte != "" {
		query = query.Where("updated_at <= ?", params.UpdatedAtLte)
	}
	if params.UpdatedAtGt != "" {
		query = query.Where("updated_at > ?", params.UpdatedAtGt)
	}
	if params.UpdatedAtLt != "" {
		query = query.Where("updated_at < ?", params.UpdatedAtLt)
	}

	if params.Page != "" && params.PageSize != "" {
		pageNum, err := strconv.Atoi(params.Page)
		if err != nil {
			pageNum = 1
		}
		pageSizeNum, err := strconv.Atoi(params.PageSize)
		if err != nil {
			pageSizeNum = 20
		}
		query = query.Offset((pageNum - 1) * pageSizeNum).Limit(pageSizeNum)
	}

	return query

}
