package main

import "database/sql"

type PostDBModel struct {
	DB *sql.DB
}

func (m *PostDBModel) CreatePost(p *PostRecord) (int64, error) {
	stmtIn, err := m.DB.Prepare("INSERT INTO post (fk_user_id, content_text) VALUES (?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return -1, err
	}

	result, err := stmtIn.Exec(
		&p.UserId,
		&p.ContentText,
	)

	if err != nil {
		return -1, err
	}

	newPostId, err := result.LastInsertId()
	return newPostId, err
}

func (m *PostDBModel) UpdateImagePath(postId int, imagePath string) error {
	stmtIn, err := m.DB.Prepare("UPDATE `post` SET content_image_path = ? WHERE id = ?")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(
		&imagePath,
		&postId,
	)

	return err
}

func (m *PostDBModel) OverwritePost(postId int, p *PostRecord) error {
	// TODO
	return nil
}

func (m *PostDBModel) DeletePost(postId int) error {
	// TODO
	return nil
}

func (m *PostDBModel) CommentPost(postId int, content string) (int, error) {
	// TODO
	return -1, nil
}

func (m *PostDBModel) LikePost(postId int) error {
	// TODO
	return nil
}
