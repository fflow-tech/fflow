package constants

import "strings"

const (
	DescOrder = "desc"
	AscOrder  = "asc"

	defaultPageIndex = 1
	defaultPageSize  = 10
	defaultSortedBy  = "created_at"
	defaultOrder     = "desc"
)

// PageQuery 分页查询对象
type PageQuery struct {
	PageIndex     int  `form:"page_index" json:"page_index,omitempty"` // 分页序号
	PageSize      int  `form:"page_size" json:"page_size,omitempty"`   // 分页大小
}

// GetOffset 获取查询数据库的 offset
func (q *PageQuery) GetOffset() int {
	if q.PageIndex <= 0 {
		q.PageIndex = defaultPageIndex
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultPageSize
	}
	return (q.PageIndex - 1) * q.PageSize
}

// GetLimit 获取查询数据库的 limit
func (q *PageQuery) GetLimit() int {
	if q.PageSize <= 0 {
		q.PageSize = defaultPageSize
	}
	return q.PageSize
}

// NewPageQuery 创建一个新的查询对象
func NewPageQuery(pageIndex int, pageSize int) *PageQuery {
	if pageIndex <= 0 {
		pageIndex = defaultPageIndex
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return &PageQuery{pageIndex, pageSize}
}

// NewDefaultPageQuery 创建一个默认的查询对象
func NewDefaultPageQuery() *PageQuery {
	return NewPageQuery(defaultPageIndex, defaultPageSize)
}

// Order 排序对象
type Order struct {
	SortedBy string `form:"sorted_by,omitempty" json:"sorted_by,omitempty"`
	Order    string `form:"order,omitempty" json:"order,omitempty"`
}

// NewOrder 排序字段和顺序
func NewOrder(sortedBy string, order string) *Order {
	return &Order{SortedBy: sortedBy, Order: order}
}

// NewDefaultOrder 创建一个默认的排序字段和顺序
func NewDefaultOrder() *Order {
	return NewOrder(defaultSortedBy, defaultOrder)
}

var (
	defaultOrderSet = "id desc"
)

// OrderStr 返回排序语句
func (o *Order) OrderStr() string {
	if o == nil {
		return defaultOrderSet
	}

	if !o.IsEmptyStr() {
		return strings.Join([]string{o.SortedBy, o.Order}, " ")
	}

	return defaultOrderSet
}

// IsEmptyStr 判断排序字段是否为空，并进行补充
func (o *Order) IsEmptyStr() bool {
	if o.Order == "" {
		return true
	}

	if o.SortedBy != "" {
		return false
	}

	o.SortedBy = "desc"
	return false
}
