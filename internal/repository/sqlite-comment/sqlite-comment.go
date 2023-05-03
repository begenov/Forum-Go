package sqlitecomment

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/begenov/Forum-Go/models"
)

type CommentSqlite struct {
	db *sql.DB
}

func NewSqliteComment(db *sql.DB) *CommentSqlite {
	return &CommentSqlite{db: db}
}

func (s *CommentSqlite) CreateComment(comment models.Comment) error {
	stmt := `INSERT INTO comment (user_id, post_id, like, dislike, text, author, date) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(stmt, &comment.UserID, &comment.PostID, &comment.Likes, &comment.Dislikes, &comment.Text, &comment.Author, &comment.Date)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentSqlite) GetCommentByPostId(id int) ([]models.Comment, error) {
	comments := []models.Comment{}
	query := `SELECT id, user_id, post_id, like, dislike, text, author, date FROM comment WHERE post_id=$1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("storage: comment by id post: %w", err)
	}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Likes, &comment.Dislikes, &comment.Text, &comment.Author, &comment.Date); err != nil {
			log.Println(err.Error())
			return nil, fmt.Errorf("storage: comment by id post: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, err
}

func (s *CommentSqlite) GetCommentById(id int) (models.Comment, error) {
	var comment models.Comment
	query := `SELECT id, user_id, post_id, like, dislike, text, author, date FROM comment WHERE id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return comment, fmt.Errorf("storage: comment by id post: %w", err)
	}
	for rows.Next() {
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Likes, &comment.Dislikes, &comment.Text, &comment.Author, &comment.Date); err != nil {
			log.Println(err.Error())
			return comment, fmt.Errorf("storage: comment by id post: %w", err)
		}
	}

	return comment, err
}
