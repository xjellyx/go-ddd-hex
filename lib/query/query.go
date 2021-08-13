package query

type Meta struct {
	PageSize int `json:"pageSize" from:"pageSize"`
	PageNum  int `json:"pageNum" from:"pageNum"`
}

func (q *Meta) Offset() int {
	if q.PageNum > 0 {
		q.PageNum = (q.PageNum - 1) * q.PageSize
	}
	return q.PageNum
}

func (q *Meta) Limit() int {
	return q.PageSize
}
