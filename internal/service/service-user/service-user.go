package serviceuser

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/begenov/Forum-Go/models"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type serverProvider interface {
	CreateUser(models.User) error
	GetUserByEmail(email string) (models.User, error)
	CrateSession(user models.User) error
	GetSessionByID(id int) (models.Session, error)
	UpdateSession(session models.Session) error
	DeleteSession(id int) error
	SessionByToken(token string) (models.Session, error)
	UserByID(id int) (models.User, error)
	DeleteToken(value string) error
}

type ServerUser struct {
	serverProvider serverProvider
}

func NewServiceUser(server serverProvider) *ServerUser {
	return &ServerUser{serverProvider: server}
}

func (s *ServerUser) CreateUser(user models.User) error {
	hash, err := generatePasswordHash(user.Password)
	if err != nil {
		return err
	}

	if err := checkUser(user); err != nil {
		return err
	}

	user.Password = hash

	if err := s.serverProvider.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *ServerUser) User(email string, password string) (models.User, error) {
	user, err := s.serverProvider.GetUserByEmail(email)
	if err != nil {
		return models.User{}, fmt.Errorf("can't get user: %w", err)
	}
	if err = compareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, fmt.Errorf("passwordhash and pasword not equal: %w", err)
	}
	session, err := s.serverProvider.GetSessionByID(user.ID)
	if err == sql.ErrNoRows {
		token, err := generateToken()
		if err != nil {
			return models.User{}, fmt.Errorf("can't generate token: %w)", err)
		}
		user.Token = token
		user.ExpiresAt = time.Now().Add(12 * time.Hour)
		session.UserId = user.ID
		if err = s.serverProvider.CrateSession(user); err != nil {
			return models.User{}, fmt.Errorf("can't create session: %w", err)
		}
		return user, nil

	}
	if err != nil {
		return models.User{}, fmt.Errorf("can't get session: %w", err)
	}
	token, err := generateToken()
	if err != nil {
		return models.User{}, fmt.Errorf("can't generate token: %w)", err)
	}
	session.Token = token
	session.ExpirationTime = time.Now().Add(12 * time.Hour)

	if err = s.serverProvider.UpdateSession(session); err != nil {
		return models.User{}, fmt.Errorf("can't update session: %w", err)
	}

	user.Token = session.Token
	user.ExpiresAt = session.ExpirationTime

	return user, nil
}

func (s *ServerUser) UserByToken(token string) (models.User, error) {
	session, err := s.serverProvider.SessionByToken(token)
	if err != nil {
		return models.User{}, err
	}
	user, err := s.serverProvider.UserByID(session.UserId)
	user.ExpiresAt = session.ExpirationTime
	return user, nil
}

func (s *ServerUser) DeleteSession(user models.User) error {
	return s.serverProvider.DeleteSession(user.ID)
}

func (s *ServerUser) UserByID(id int) (models.User, error) {
	return s.serverProvider.UserByID(id)
}

func (s *ServerUser) DeleteToken(value string) error {
	return s.serverProvider.DeleteToken(value)
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func compareHashAndPassword(hashpassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashpassword, password)
}

func generateToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

// passwordValidate = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
var emailValidate = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func checkUser(user models.User) error {
	if !emailValidate.MatchString(user.Email) {
		return fmt.Errorf("Error: email incorect")
	}

	if len(user.Name) <= 2 || len(user.Name) >= 29 {
		return fmt.Errorf("username incorect")
	}

	if len(user.Password) <= 7 {
		return fmt.Errorf("password incorect")
	}

	return nil
}
