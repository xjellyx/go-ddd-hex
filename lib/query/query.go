package query

import "gorm.io/gorm"

type Meta struct {
	PageSize int `json:"pageSize" from:"pageSize"`
	PageNum  int `json:"pageNum" from:"pageNum"`
}

func (q *Meta) WithOffsetLimit(db *gorm.DB) *Meta {
	if q.PageNum > 0 {
		q.PageNum = (q.PageNum - 1) * q.PageSize
		db.Offset(q.PageNum)
	}

	if q.PageSize > 0 {
		db.Limit(q.PageSize)
	}
	return q
}
