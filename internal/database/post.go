package database

import (
	"database/sql"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
)

type PostDB interface {
	// New create new post and return latest id
	New(payload *payloads.PostPayload) (int64, error)

	// Edit overwrite post for editting
	Edit(id int, payload *payloads.PostPayload) error

	// Delete soft delete for post
	Delete(id int) error

	// Comment add comment to post and return new cmt id
	Comment(postId int, userId int, content string) (int64, error)

	// Like add like to post
	Like(postId int, userId int) error

	// GetPostById load post given id
	GetPostById(id int) (payloads.PostPayload, error)

	// ListPostByUserId load post from given user-id
	ListPostByUserId(userId int) ([]payloads.PostPayload, error)

	// ListPostByUserIdAndFollower return all post of current user-id and follower
	ListPostByUserIdAndFollower(userId int) ([]payloads.PostPayload, error)

	// UpdateImagePath overwrite image path for given post
	UpdateImagePath(postId int, imagePath string) error
}

func NewPostDB(db *sql.DB) PostDB {
	return &postDB{db: db}
}

type postDB struct {
	db *sql.DB
}

func (m *postDB) ListPostByUserIdAndFollower(userId int) ([]payloads.PostPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (m *postDB) GetPostById(id int) (payloads.PostPayload, error) {
	stmtOut, err := m.db.Prepare("SELECT id, fk_user_id, content_text, IFNULL(content_image_path, '') AS content_image_path, created_at FROM `post` WHERE id = ? AND visible = 1")
	defer stmtOut.Close()

	var p payloads.PostPayload

	if err != nil {
		return p, err
	}

	err = stmtOut.QueryRow(id).Scan(
		&p.Id, &p.UserId, &p.ContentText, &p.ContentImagePath,
		&p.CreatedAt,
	)

	return p, err
}

func (m *postDB) New(payload *payloads.PostPayload) (int64, error) {
	stmtIn, err := m.db.Prepare("INSERT INTO post (fk_user_id, content_text) VALUES (?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return -1, err
	}

	result, err := stmtIn.Exec(
		&payload.UserId,
		&payload.ContentText,
	)

	if err != nil {
		return -1, err
	}

	newPostId, err := result.LastInsertId()
	return newPostId, err
}

func (m *postDB) UpdateImagePath(postId int, imagePath string) error {
	stmtIn, err := m.db.Prepare("UPDATE `post` SET content_image_path = ? WHERE id = ?")
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

func (m *postDB) Edit(postId int, payload *payloads.PostPayload) error {
	stmtIn, err := m.db.Prepare("UPDATE `post` SET content_text = ? WHERE id = ?")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(
		&payload.ContentText,
		&postId,
	)

	return err
}

func (m *postDB) Delete(id int) error {
	stmtIn, err := m.db.Prepare("UPDATE `post` SET visible = 0 WHERE id = ?")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(&id)

	return err

}

func (m *postDB) Comment(postId int, userId int, content string) (int64, error) {
	stmtIn, err := m.db.Prepare("INSERT INTO comment (fk_post_id, fk_user_id, content) VALUES (?, ?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return -1, err
	}

	result, err := stmtIn.Exec(
		&postId,
		&userId,
		&content,
	)

	if err != nil {
		return -1, err
	}

	newCmtId, err := result.LastInsertId()
	return newCmtId, err
}

func (m *postDB) Like(postId int, userId int) error {
	stmtIn, err := m.db.Prepare("INSERT INTO `like` (fk_post_id, fk_user_id) VALUES (?, ?)")
	defer stmtIn.Close()

	if err != nil {
		return err
	}

	_, err = stmtIn.Exec(
		&postId,
		&userId,
	)

	return err
}

func (m *postDB) ListPostByUserId(userId int) ([]payloads.PostPayload, error) {
	var posts []payloads.PostPayload

	stmtOut, err := m.db.Prepare("SELECT id, content_text, IFNULL(content_image_path, '') AS content_image_path, created_at FROM `user_user` u_u LEFT JOIN `post` p ON u_u.fk_follower_id = p.fk_user_id WHERE u_u.fk_user_id = ? AND visible = 1")
	defer stmtOut.Close()

	if err != nil {
		return posts, err
	}

	rows, err := stmtOut.Query(userId)
	defer rows.Close()

	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var p payloads.PostPayload

		err := rows.Scan(&p.Id, &p.ContentText, &p.ContentImagePath, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	err = rows.Err()
	if err != nil {
		return posts, err
	}

	return posts, nil
}

//func (m UserDBModel) ViewPosts(id int) ([]models2.PostRecord, error) {
//	stmtOut, err := m.DB.Prepare("SELECT id, content_text, IFNULL(content_image_path, '') AS content_image_path, created_at FROM `post` p WHERE p.fk_user_id = ? AND visible = 1")
//	defer stmtOut.Close()
//
//	if err != nil {
//		return nil, err
//	}
//
//	rows, err := stmtOut.Query(id)
//	defer rows.Close()
//
//	if err != nil {
//		return nil, err
//	}
//
//	err = rows.Err()
//	if err != nil {
//		return nil, err
//	}
//
//	var posts []models2.PostRecord
//	for rows.Next() {
//		var p models2.PostRecord
//
//		err := rows.Scan(&p.Id, &p.ContentText, &p.ContentImagePath, &p.CreatedAt)
//		if err != nil {
//			return nil, err
//		}
//
//		posts = append(posts, p)
//	}
//
//	return posts, nil
//}
