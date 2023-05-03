package controllerreaction

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/begenov/Forum-Go/models"
	"github.com/begenov/Forum-Go/pkg"
)

type reactionProvider interface {
	ReactionPost(post_reaction *models.ReactionPost) error
	ReactionComment(comment_reaction *models.ReactionComment) error
}

type commentProvider interface {
	GetCommentById(id int) (models.Comment, error)
}

type userProvider interface {
	UserByID(int) (models.User, error)
}

type postProvider interface {
	GetPostByID(id int) (models.Post, error)
}

type ReactionController struct {
	reaction reactionProvider
	user     userProvider
	comment  commentProvider
	post     postProvider
}

func NewReactionController(reaction reactionProvider, userreactionProvider userProvider, comment commentProvider, post postProvider) *ReactionController {
	return &ReactionController{
		reaction: reaction,
		user:     userreactionProvider,
		comment:  comment,
		post:     post,
	}
}

func (c *ReactionController) ReactionPost(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.UserIDKey)

	if userID.(int) <= 0 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	user, err := c.user.UserByID(userID.(int))
	if err != nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if postId <= 0 || err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	post, err := c.post.GetPostByID(postId)
	if err != nil || post.Id == 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res := r.Form.Get("isLike")
	if strings.TrimSpace(res) == "" {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if res == "like" {
		err := c.reaction.ReactionPost(&models.ReactionPost{PostID: post.Id, UserID: user.ID, Islike: 1})
		if err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if res == "dislike" {
		err := c.reaction.ReactionPost(&models.ReactionPost{PostID: post.Id, UserID: user.ID, Islike: 0})
		if err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	page := fmt.Sprintf("/post?id=%d", post.Id)
	http.Redirect(w, r, page, http.StatusSeeOther)
}

func (c *ReactionController) ReactionComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.UserIDKey)

	if userID.(int) <= 0 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	user, err := c.user.UserByID(userID.(int))
	if err != nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	commentId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	comment, err := c.comment.GetCommentById(commentId)
	if err != nil || comment.ID == 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if postId <= 0 || err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	post, err := c.post.GetPostByID(postId)
	if post.Id == 0 || err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res := r.Form.Get("islike")
	if res == "" {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if res == "like" {
		err = c.reaction.ReactionComment(&models.ReactionComment{UserID: user.ID, CommentID: comment.ID, Islike: 1})
		if err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else if res == "dislike" {
		err = c.reaction.ReactionComment(&models.ReactionComment{UserID: user.ID, CommentID: comment.ID, Islike: 0})
		if err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	page := fmt.Sprintf("/post?id=%d", postId)
	http.Redirect(w, r, page, http.StatusSeeOther)
	return
}
