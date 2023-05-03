package commentreaction

import (
	"database/sql"
	"log"

	"github.com/begenov/Forum-Go/models"
)

type CommentReaction struct {
	db *sql.DB
}

type count struct {
	likes    int
	dislikes int
}

func NewCommentReaction(db *sql.DB) *CommentReaction {
	return &CommentReaction{db: db}
}

func (c *CommentReaction) CreateCommentReaction(comment models.ReactionComment) error {
	stmt := `INSERT INTO reaction_comment (user_id, comment_id, like_is) VALUES (?, ?, ?)`
	if _, err := c.db.Exec(stmt, comment.UserID, comment.CommentID, comment.Islike); err != nil {
		return err
	}
	return c.UpdateOnComment(comment.CommentID)
}

func (c *CommentReaction) UpdateOnComment(id int) error {
	var count count
	query := `SELECT
	(SELECT COUNT(*) FROM reaction_comment WHERE comment_id = $1 AND like_is = 1) AS likes,
	(SELECT COUNT(*) FROM reaction_comment WHERE comment_id = $1 AND like_is = 0) AS dislikes`
	row := c.db.QueryRow(query, id)
	if err := row.Scan(&count.likes, &count.dislikes); err != nil {
		return err
	}
	query = `UPDATE comment SET like = ?, dislike = ? WHERE id = ?`
	if _, err := c.db.Exec(query, count.likes, count.dislikes, id); err != nil {
		return err
	}
	return nil
}

func (r *CommentReaction) CheckCommentReaction(commentReac *models.ReactionComment) (bool, error) {
	var exists bool
	{
		query := `SELECT EXISTS(SELECT 1 FROM reaction_comment WHERE comment_id = ? AND user_id = ?) AS value_exists;`
		err := r.db.QueryRow(query, &commentReac.CommentID, &commentReac.UserID).Scan(&exists)
		if err != nil {
			log.Println(err.Error(), " 50")
			return exists, err
		}
	}
	{
		var exists1 bool
		query := `SELECT EXISTS(SELECT 1 FROM reaction_comment WHERE comment_id = $1  AND user_id = $2 AND like_is = $3) value_exists;`

		err := r.db.QueryRow(query, commentReac.CommentID, commentReac.UserID, commentReac.Islike).Scan(&exists1)
		if err != nil {
			log.Println(err.Error(), " 60")
			return exists, err
		}
		if exists1 {
			commentReac.Islike = -1
		}
	}
	return exists, nil
}

func (r *CommentReaction) UpdateReactionComment(commentReac models.ReactionComment) error {
	query := `UPDATE reaction_comment SET like_is = $1 WHERE user_id = $2 AND comment_id = $3`
	if _, err := r.db.Exec(query, commentReac.Islike, commentReac.UserID, commentReac.CommentID); err != nil {
		log.Println(err.Error(), " 74")
		return err
	}
	return r.UpdateOnComment(commentReac.CommentID)
}
