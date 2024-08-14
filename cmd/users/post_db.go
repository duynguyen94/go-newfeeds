package main

import "database/sql"

type PostDBModel struct {
	DB *sql.DB
}

func (m *PostDBModel) CreatePost(p *PostRecord) (int, error) {
	// TODO
	return -1, nil
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
