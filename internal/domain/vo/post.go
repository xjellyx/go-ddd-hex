package vo

type PostRes struct {
	BaseRes
	Title   string `json:"title"`
	Content string `json:"content"`
}
