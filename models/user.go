package models

import "time"

type contextKey string

const UserIDKey contextKey = "userID"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Token     string
	ExpiresAt time.Time
}

type Session struct {
	ID             int
	UserId         int
	Token          string
	ExpirationTime time.Time
}

type TemplateData struct {
	UserID          int
	User            string
	Post            Post
	Posts           []Post
	Comment         []Comment
	IsAuthenticated bool
	MsgError        string
}

type Comment struct {
	ID       int
	UserID   int
	PostID   int
	Likes    int
	Dislikes int
	Text     string
	Author   string
	Date     string
	IsLike   int
}

type Post struct {
	Id          int
	UserID      int
	Title       string
	Description string
	Category    []string
	Author      string
	Likes       int
	Dislikes    int
	IsLike      int
	CreateAt    string
}

type ReactionComment struct {
	ID        int
	UserID    int
	CommentID int
	Islike    int
}

type ReactionPost struct {
	ID     int
	PostID int
	UserID int
	Islike int
}
