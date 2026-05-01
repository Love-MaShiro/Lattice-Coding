package db

type PageQuery struct {
	Page    int    `form:"page" json:"page"`
	Size    int    `form:"size" json:"size"`
	Keyword string `form:"keyword" json:"keyword"`
}

func (p *PageQuery) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 20
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return (p.Page - 1) * p.Size
}

func (p *PageQuery) Limit() int {
	if p.Size <= 0 {
		p.Size = 20
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return p.Size
}

type PageResult[T any] struct {
	Items     []T  `json:"items"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalPage int   `json:"total_page"`
}

func NewPageResult[T any](query *PageQuery, items []T, total int64) PageResult[T] {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	size := query.Limit()
	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}
	return PageResult[T]{
		Items:     items,
		Total:     total,
		Page:      page,
		Size:      size,
		TotalPage: totalPage,
	}
}
