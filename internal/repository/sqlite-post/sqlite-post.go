package sqlitepost

import (
	"database/sql"
	"log"
	"strings"

	"github.com/begenov/Forum-Go/models"
)

type SqlitePost struct {
	db *sql.DB
}

func NewSqlitePost(db *sql.DB) *SqlitePost {
	return &SqlitePost{
		db: db,
	}
}

func (s *SqlitePost) CreatePost(post models.Post) (int, error) {
	category := strings.Join(post.Category, ", ")

	stmt := `INSERT INTO post (author_id, like, dislike, title, category, content, author, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := s.db.Exec(stmt, &post.UserID, &post.Likes, &post.Dislikes, &post.Title, &category, &post.Description, &post.Author, &post.CreateAt)
	if err != nil {
		log.Printf("error post create: %v", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = s.createCategory(post.Category, int(id))

	return int(id), nil
}

func (s *SqlitePost) GetPostByID(id int) (models.Post, error) {
	query := `SELECT id, author_id, title, like, dislike, title, category, content, author, date FROM post WHERE id = ?`

	row := s.db.QueryRow(query, id)
	var post models.Post
	var category string
	if err := row.Scan(&post.Id, &post.UserID, &post.Title, &post.Likes, &post.Dislikes, &post.Title, &category, &post.Description, &post.Author, &post.CreateAt); err != nil {
		return post, err
	}
	post.Category = strings.Split(category, " ")
	return post, nil
}

func (s *SqlitePost) GetAllPost() ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post`
	row, err := s.db.Query(query)
	if err != nil {
		log.Printf("error get all post exec %v", err)
		return posts, err
	}
	for row.Next() {
		var post models.Post
		var category string
		if err := row.Scan(&post.Id, &post.UserID, &post.Likes, &post.Dislikes, &post.Title, &category, &post.Description, &post.Author, &post.CreateAt); err != nil {
			log.Printf("error get all post scan %v", err)
			return posts, err
		}
		post.Category = append(post.Category, category)
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *SqlitePost) GetPostsByUserID(userID int) ([]models.Post, error) {
	query := `SELECT * FROM post WHERE author_id=?`
	var posts []models.Post
	row, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var category string
		var post models.Post
		if err := row.Scan(&post.Id, &post.UserID, &post.Likes, &post.Dislikes, &post.Title, &category, &post.Description, &post.Author, &post.CreateAt); err != nil {
			return nil, err
		}
		post.Category = strings.Split(category, " ")
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *SqlitePost) GetMyLikedPosts(author_id int) ([]models.Post, error) {
	var posts []models.Post

	query := `
		SELECT 
			p.id,
			p.author_id,
			p.like,
			p.dislike,
			p.title,
			p.category,
			p.content,
			p.author,
			p.date
		FROM 
			post p
		JOIN
			reaction_post rp
		ON
			p.id = rp.post_id
		WHERE 
			rp.user_id = ? 
		AND 
			rp.like_is = 1`
	/*
					JOIN
		        likesPost lp
		    ON
		        p.postId = lp.postId
		    WHERE
		        lp.userId = $1 AND lp.like1 = 1
	*/
	row, err := s.db.Query(query, author_id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var post models.Post
		var category string
		if err := row.Scan(&post.Id, &post.UserID, &post.Likes, &post.Dislikes, &post.Title, &category, &post.Description, &post.Author, &post.CreateAt); err != nil {
			return nil, err
		}
		post.Category = strings.Split(category, " ")
		posts = append(posts, post)

	}
	return posts, nil
}

func (s *SqlitePost) createCategory(category []string, postID int) error {
	query := `INSERT INTO category(post_id, category) VALUES (?, ?);`
	for _, v := range category {
		_, err := s.db.Exec(query, postID, v)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *SqlitePost) GetPostByCategory(category []string) ([]models.Post, error) {
	var posts []models.Post

	slice, err := s.getCategory(category)
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM post WHERE id = ?`
	for _, v := range slice {
		row, err := s.db.Query(query, v)
		if err != nil {
			return nil, err
		}
		for row.Next() {
			var post models.Post
			var category string
			if err := row.Scan(&post.Id, &post.UserID, &post.Likes, &post.Dislikes, &post.Title, &category, &post.Description, &post.Author, &post.CreateAt); err != nil {
				return nil, err
			}
			if isId(post, posts) {
				post.Category = strings.Split(category, " ")
				posts = append(posts, post)
			}

		}
	}

	return posts, nil
}

func (s *SqlitePost) getCategory(category []string) ([]int, error) {
	var slice []int

	query := `SELECT post_id FROM category WHERE category = ?`
	for _, v := range category {
		row, err := s.db.Query(query, v)
		if err != nil {
			return nil, err
		}
		for row.Next() {
			var i int
			if err := row.Scan(&i); err != nil {
				return nil, err
			}
			slice = append(slice, i)
		}
	}
	return slice, nil
}

func isId(post models.Post, posts []models.Post) bool {
	for _, v := range posts {
		if v.Id == post.Id {
			return false
		}
	}
	return true
}
