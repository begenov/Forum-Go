package sqliteuser

import (
	"database/sql"
	"fmt"

	"github.com/begenov/Forum-Go/models"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite(db *sql.DB) *UserSqlite {
	return &UserSqlite{db: db}
}

func (s *UserSqlite) CreateUser(user models.User) error {
	stmt := `INSERT INTO user (email, username, password_hash) VALUES (?, ?, ?)`
	_, err := s.db.Exec(stmt, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserSqlite) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	stmt := `SELECT id, email, username, password_hash FROM user WHERE email = ?`
	row := s.db.QueryRow(stmt, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserSqlite) UserByID(id int) (models.User, error) {
	var user models.User
	stmt := `SELECT id, email, username, password_hash FROM user WHERE id = ?`
	row := s.db.QueryRow(stmt, id)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserSqlite) CrateSession(user models.User) error {
	stmt := `INSERT INTO session (user_id, token, expiration_time) VALUES (?, ?, ?)`

	if _, err := s.db.Exec(stmt, &user.ID, &user.Token, &user.ExpiresAt); err != nil {
		return fmt.Errorf("can't create session: %w", err)
	}

	return nil
}

func (s *UserSqlite) GetSessionByID(id int) (models.Session, error) {
	var session models.Session
	stmt := `SELECT * FROM session WHERE user_id = ?`
	if err := s.db.QueryRow(stmt, id).Scan(&session.ID, &session.UserId, &session.Token, &session.ExpirationTime); err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (s *UserSqlite) UpdateSession(session models.Session) error {
	stmt := `UPDATE session SET token = ?, expiration_time = ? WHERE user_id = ?`

	if _, err := s.db.Exec(stmt, &session.Token, &session.ExpirationTime, &session.UserId); err != nil {
		return fmt.Errorf("can't update user session: %w", err)
	}

	return nil
}

func (s *UserSqlite) SessionByToken(token string) (models.Session, error) {
	var session models.Session
	query := `SELECT id, user_id, token, expiration_time FROM session WHERE token = ?`
	row := s.db.QueryRow(query, token)
	err := row.Scan(&session.ID, &session.UserId, &session.Token, &session.ExpirationTime)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (s *UserSqlite) DeleteToken(value string) error {
	stmt := `DELETE FROM session WHERE token = ?`
	_, err := s.db.Exec(stmt, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserSqlite) DeleteSession(id int) error {
	stmt := `DELETE FROM session WHERE user_id = ?`
	_, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
