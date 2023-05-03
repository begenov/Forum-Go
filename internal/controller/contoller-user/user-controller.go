package contolleruser

import (
	"net/http"
	"text/template"

	"github.com/begenov/Forum-Go/models"
	"github.com/begenov/Forum-Go/pkg"
)

type userProvider interface {
	CreateUser(user models.User) error
	User(email string, password string) (models.User, error)
	UserByID(int) (models.User, error)
	DeleteToken(value string) error
}

type UserController struct {
	user userProvider
}

func NewControllerUser(user userProvider) *UserController {
	return &UserController{user: user}
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var templateData models.TemplateData
	userID := r.Context().Value(models.UserIDKey)
	if userID.(int) > 0 {
		user, _ := c.user.UserByID(userID.(int))
		if user != (models.User{}) {
			pkg.ResponseError(w, "Error to start with, do logout", http.StatusBadRequest)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("./templates/html/sign-up.html")
		if err != nil {
			pkg.ResponseError(w, "Status Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			pkg.ResponseError(w, "Status Internal Server Error", http.StatusInternalServerError)
			return
		}
		user := getUser(r)
		if err := c.user.CreateUser(user); err != nil {
			templateData.MsgError = err.Error()
			tmpl, err := template.ParseFiles("./templates/html/sign-up.html")
			if err != nil {
				pkg.ResponseError(w, "Status Internal Server Error", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, templateData)
			return
		}
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
	default:
		pkg.ResponseError(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (c *UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	var templateData models.TemplateData
	userID := r.Context().Value(models.UserIDKey)
	if userID.(int) > 0 {
		user, _ := c.user.UserByID(userID.(int))
		if user != (models.User{}) {
			pkg.ResponseError(w, "Error to start with, do logout", http.StatusBadRequest)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("./templates/html/sign-in.html")
		if err != nil {
			pkg.ResponseError(w, "Status Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		user := getUser(r)
		user, err := c.user.User(user.Email, user.Password)
		if err != nil {
			templateData.MsgError = "Incorect: email or passoword"
			tmpl, err := template.ParseFiles("./templates/html/sign-in.html")
			if err != nil {
				pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, templateData)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    user.Token,
			Expires:  user.ExpiresAt,
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (c *UserController) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	if err := c.user.DeleteToken(cookie.Value); err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie.MaxAge = -1
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUser(r *http.Request) models.User {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	return models.User{
		Email:    email,
		Name:     username,
		Password: password,
	}
}
