package servicereaction

import "github.com/begenov/Forum-Go/models"

type (
	postProvider interface {
		CreateReactionPost(reaction_post models.ReactionPost) error
		CheckReactionPost(post_reaction *models.ReactionPost) (bool, error)
		UpdateReactionPost(post_reaction models.ReactionPost) error
	}
	commentProvider interface {
		CreateCommentReaction(comment models.ReactionComment) error
		UpdateReactionComment(commentReac models.ReactionComment) error
		CheckCommentReaction(commentReac *models.ReactionComment) (bool, error)
	}
)

type ReactionService struct {
	post    postProvider
	comment commentProvider
}

func NewServiceReaction(post postProvider, comment commentProvider) *ReactionService {
	return &ReactionService{
		post:    post,
		comment: comment,
	}
}

func (r *ReactionService) ReactionPost(post_reaction *models.ReactionPost) error {
	isExist, err := r.post.CheckReactionPost(post_reaction)
	if err != nil {
		return err
	}
	if !isExist {
		return r.post.CreateReactionPost(*post_reaction)
	}
	return r.post.UpdateReactionPost(*post_reaction)
}

func (r *ReactionService) ReactionComment(comment_reaction *models.ReactionComment) error {
	isExist, err := r.comment.CheckCommentReaction(comment_reaction)
	if err != nil {
		return err
	}
	if !isExist {
		return r.comment.CreateCommentReaction(*comment_reaction)
	}
	return r.comment.UpdateReactionComment(*comment_reaction)
}
