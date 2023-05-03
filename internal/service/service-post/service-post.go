package servicepost

import (
	"fmt"
	"strings"

	"github.com/begenov/Forum-Go/models"
)

type postProvider interface {
	CreatePost(post models.Post) (int, error)
	GetPostByID(id int) (models.Post, error)
	GetAllPost() ([]models.Post, error)
	GetPostsByUserID(userID int) ([]models.Post, error)
	GetMyLikedPosts(author_id int) ([]models.Post, error)
	GetPostByCategory(category []string) ([]models.Post, error)
}

type ServicePost struct {
	postProvider postProvider
}

func NewServicePost(postProvider postProvider) *ServicePost {
	return &ServicePost{postProvider: postProvider}
}

func (s *ServicePost) CreatePost(post models.Post) (int, error) {
	if err := checkPost(post); err != nil {
		return 0, err
	}
	return s.postProvider.CreatePost(post)
}

func (s *ServicePost) GetPostByID(id int) (models.Post, error) {
	return s.postProvider.GetPostByID(id)
}

func (s *ServicePost) AllPost() ([]models.Post, error) {
	return s.postProvider.GetAllPost()
}

func (s *ServicePost) GetPostsByUserID(userID int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("error invalid argument")
	}
	return s.postProvider.GetPostsByUserID(userID)
}

func (s *ServicePost) GetMyLikedPosts(userID int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("error invalid argument")
	}
	return s.postProvider.GetMyLikedPosts(userID)
}

func (s *ServicePost) GetPostsByCategory(category []string) ([]models.Post, error) {
	if err := categoryCheck(category); err != nil {
		return nil, err
	}
	return s.postProvider.GetPostByCategory(category)
}

// Checker Post

func checkPost(post models.Post) error {
	if strings.TrimSpace(post.Title) == "" {
		return fmt.Errorf("empty Title")
	}

	if err := categoryCheck(post.Category); err != nil {
		return err
	}

	if strings.TrimSpace(post.Description) == "" {
		return fmt.Errorf("empty Content")
	}
	return nil
}

func categoryCheck(category1 []string) error {
	if len(category1) == 0 {
		return fmt.Errorf("empty category")
	}

	category := []string{"Technology", "Science", "Art", "Sports", "Music"}
	for _, v := range category {
		for _, j := range category1 {
			if v == j {
				return nil
			}
		}
	}
	return fmt.Errorf("you cannot select another category")
}
