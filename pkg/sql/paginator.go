package sql

import (
	"context"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"github.com/samber/lo"
	"github.com/ywengineer/smart/utility"
	"gorm.io/gorm"
	"math"
)

func CreatePaginator(pageSize int, rules []paginator.Rule) *paginator.Paginator {
	return paginator.New(
		&paginator.Config{
			Rules: rules,
			Limit: pageSize,
		},
	)
}

type PageInfo[T interface{}] struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
	Total    int64 `json:"total"`
	MaxPage  int32 `json:"max_page"`
	HasNext  bool  `json:"has_next"`
	Data     []T   `json:"data"`
}

// Paginate 分页工具函数
func Paginate[T interface{}, V interface{}](ctx context.Context, query *gorm.DB, page, pageSize int32, dest *[]T, converter func(i T) V) (*PageInfo[V], error) {
	var total int64
	offset := (utility.MaxInt(page, 1) - 1) * pageSize
	// 查询总数
	if err := query.WithContext(ctx).Count(&total).Error; err != nil {
		return nil, err
	}
	//
	maxPage := int32(math.Ceil(float64(total) / float64(pageSize)))
	if page > maxPage {
		return &PageInfo[V]{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			MaxPage:  maxPage,
			HasNext:  false,
			Data:     []V{},
		}, nil
	}
	// 分页查询
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(dest).Error; err != nil {
		return nil, err
	}
	var ret []V
	if converter != nil {
		ret = lo.Map(*dest, func(item T, index int) V {
			return converter(item)
		})
	}
	return &PageInfo[V]{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		MaxPage:  maxPage,
		HasNext:  page < maxPage,
		Data:     ret,
	}, nil
}
