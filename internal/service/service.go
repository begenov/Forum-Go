package service

import (
	"github.com/begenov/Forum-Go/internal/repository"
	servicecomment "github.com/begenov/Forum-Go/internal/service/service-comment"
	servicepost "github.com/begenov/Forum-Go/internal/service/service-post"
	servicereaction "github.com/begenov/Forum-Go/internal/service/service-reaction"
	serviceuser "github.com/begenov/Forum-Go/internal/service/service-user"
)

type Service struct {
	User     serviceuser.ServerUser
	Post     servicepost.ServicePost
	Comment  servicecomment.CommentService
	Reaction servicereaction.ReactionService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:     *serviceuser.NewServiceUser(&repos.User),
		Post:     *servicepost.NewServicePost(&repos.Post),
		Comment:  *servicecomment.NewCommentService(&repos.Comment),
		Reaction: *servicereaction.NewServiceReaction(repos.Reaction.PostReaction, repos.Reaction.CommentReaction),
	}
}
