package servicecomment

import (
	"fmt"
	"strings"

	"github.com/begenov/Forum-Go/models"
)

type commentProvider interface {
	CreateComment(comment models.Comment) error
	GetCommentByPostId(id int) ([]models.Comment, error)
	GetCommentById(id int) (models.Comment, error)
}

type CommentService struct {
	comment commentProvider
}

func NewCommentService(commentProvider commentProvider) *CommentService {
	return &CommentService{comment: commentProvider}
}

func (s *CommentService) CreateComment(c models.Comment) error {
	if strings.TrimSpace(c.Text) == "" {
		return fmt.Errorf("Error: empty text")
	}
	return s.comment.CreateComment(c)
}

func (s *CommentService) GetCommentByPostId(id int) ([]models.Comment, error) {
	return s.comment.GetCommentByPostId(id)
}

func (s *CommentService) GetCommentById(id int) (models.Comment, error) {
	return s.comment.GetCommentById(id)
}
