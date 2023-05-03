package postreaction

import (
	"database/sql"
	"log"

	"github.com/begenov/Forum-Go/models"
)

type PostReaction struct {
	db *sql.DB
}

type count struct {
	likes    int
	dislikes int
}

func NewReactionPost(db *sql.DB) *PostReaction {
	return &PostReaction{db: db}
}

func (p *PostReaction) CreateReactionPost(reaction_post models.ReactionPost) error {
	query := `INSERT INTO reaction_post(user_id, post_id, like_is) VALUES ($1, $2, $3)`
	if _, err := p.db.Exec(query, &reaction_post.UserID, &reaction_post.PostID, &reaction_post.Islike); err != nil {
		return err
	}
	return p.UpdateOnPost(reaction_post.PostID)
}

func (p *PostReaction) UpdateOnPost(id int) error {
	var count count
	query := `SELECT
	(SELECT COUNT(*) FROM reaction_post WHERE post_id = $1 AND like_is = 1) AS likes,
	(SELECT COUNT(*) FROM reaction_post WHERE post_id = $1 AND like_is = 0) AS dislike;`

	row := p.db.QueryRow(query, id)
	if err := row.Scan(&count.likes, &count.dislikes); err != nil {
		log.Println("okkk")
		return err
	}
	query = `UPDATE post SET like = $1, dislike = $2 WHERE id = $3`

	if _, err := p.db.Exec(query, count.likes, count.dislikes, id); err != nil {
		return err
	}
	return nil
}

func (p *PostReaction) CheckReactionPost(post_reaction *models.ReactionPost) (bool, error) {
	var exists bool
	{
		query := `SELECT EXISTS(SELECT 1 FROM reaction_post WHERE post_id = $1 and user_id = $2) AS value_exists;`
		err := p.db.QueryRow(query, post_reaction.PostID, post_reaction.UserID).Scan(&exists)
		if err != nil {
			return exists, err
		}
	}

	{
		var exists1 bool
		query := `SELECT EXISTS(SELECT 1 FROM reaction_post WHERE post_id = $1 and user_id = $2 and like_is = $3) AS value_exists;`
		err := p.db.QueryRow(query, post_reaction.PostID, post_reaction.UserID, post_reaction.Islike).Scan(&exists1)
		if err != nil {
			return exists, err
		}
		if exists1 {
			post_reaction.Islike = -1
		}
	}
	return exists, nil
}

func (p *PostReaction) UpdateReactionPost(post_reaction models.ReactionPost) error {
	query := `UPDATE reaction_post SET like_is = $1 WHERE user_id = $2 AND post_id = $3`
	_, err := p.db.Exec(query, post_reaction.Islike, post_reaction.UserID, post_reaction.PostID)
	if err != nil {
		return err
	}
	return p.UpdateOnPost(post_reaction.PostID)
}
