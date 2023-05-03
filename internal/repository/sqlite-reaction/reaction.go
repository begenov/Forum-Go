package sqlitereaction

import (
	"database/sql"

	commentreaction "github.com/begenov/Forum-Go/internal/repository/sqlite-reaction/comment-reaction"
	postreaction "github.com/begenov/Forum-Go/internal/repository/sqlite-reaction/post-reaction"
	"github.com/begenov/Forum-Go/models"
)

type PostReaction interface {
	CreateReactionPost(reaction_post models.ReactionPost) error
	CheckReactionPost(post_reaction *models.ReactionPost) (bool, error)
	UpdateReactionPost(post_reaction models.ReactionPost) error
}

type CommentReaction interface {
	CreateCommentReaction(comment models.ReactionComment) error
	UpdateReactionComment(commentReac models.ReactionComment) error
	CheckCommentReaction(commentReac *models.ReactionComment) (bool, error)
}

type Reaction struct {
	PostReaction    PostReaction
	CommentReaction CommentReaction
}

func NewReaction(db *sql.DB) *Reaction {
	return &Reaction{
		PostReaction:    postreaction.NewReactionPost(db),
		CommentReaction: commentreaction.NewCommentReaction(db),
	}
}
