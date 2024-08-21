package payloads

type CommentRecord struct {
	ContentText string `json:"content"`
	UserId      int    `json:"userId"`
}
