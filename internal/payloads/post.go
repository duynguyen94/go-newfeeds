package payloads

type PostPayload struct {
	Id               int    `json:"id,omitempty"`
	ContentText      string `json:"text,omitempty"`
	ContentImagePath string `json:"imagePath"`
	CreatedAt        string `json:"createdAt,omitempty"`
	UserId           int    `json:"userId,omitempty"`
	DownloadUrl      string `json:"downloadUrl,omitempty"`
}

func (p *PostPayload) Merge(newPost *PostPayload) {
	if newPost.ContentText != "" {
		p.ContentText = newPost.ContentText
	}

	// What about image???
}
