package vo

import (
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"strconv"
)

type PostVO struct {
	BaseVO
	Title   string `json:"title"`
	Content string `json:"content"`
}

func PostEntity2VO(in *entity.Post) *PostVO {
	res := &PostVO{
		Title:   in.Title,
		Content: in.Content,
	}

	res.CreatedAt = in.CreatedAt
	res.UpdatedAt = in.UpdatedAt
	res.ID = strconv.Itoa(int(in.ID))

	return res
}
