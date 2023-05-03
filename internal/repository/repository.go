package repository

import (
	"database/sql"

	sqlitecomment "github.com/begenov/Forum-Go/internal/repository/sqlite-comment"
	sqlitepost "github.com/begenov/Forum-Go/internal/repository/sqlite-post"
	sqlitereaction "github.com/begenov/Forum-Go/internal/repository/sqlite-reaction"
	sqliteuser "github.com/begenov/Forum-Go/internal/repository/sqlite-user"
)

type Repository struct {
	User     sqliteuser.UserSqlite
	Post     sqlitepost.SqlitePost
	Comment  sqlitecomment.CommentSqlite
	Reaction sqlitereaction.Reaction
}

func NewRepositroy(db *sql.DB) *Repository {
	return &Repository{
		User:     *sqliteuser.NewUserSqlite(db),
		Post:     *sqlitepost.NewSqlitePost(db),
		Comment:  *sqlitecomment.NewSqliteComment(db),
		Reaction: *sqlitereaction.NewReaction(db),
	}
}
