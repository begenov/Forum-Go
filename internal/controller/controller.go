package controller

import (
	"html/template"
	"net/http"

	contolleruser "github.com/begenov/Forum-Go/internal/controller/contoller-user"
	controllerpost "github.com/begenov/Forum-Go/internal/controller/controller-post"
	controllerreaction "github.com/begenov/Forum-Go/internal/controller/controller-reaction"
	"github.com/begenov/Forum-Go/internal/service"
	"github.com/begenov/Forum-Go/models"
	"github.com/begenov/Forum-Go/pkg"
)

type Controller struct {
	User     *contolleruser.UserController
	Post     *controllerpost.PostController
	Reaction *controllerreaction.ReactionController
	service  service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{
		User:     contolleruser.NewControllerUser(&service.User),
		Post:     controllerpost.NewPostController(&service.Post, &service.User, &service.Comment),
		Reaction: controllerreaction.NewReactionController(&service.Reaction, &service.User, &service.Comment, &service.Post),
		service:  service,
	}
}

func (c *Controller) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		pkg.ResponseError(w, "Status Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		pkg.ResponseError(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var tempalteDate models.TemplateData
	userID := r.Context().Value(models.UserIDKey)

	if userID.(int) > 0 {
		user, err := c.service.User.UserByID(userID.(int))
		if err == nil {
			tempalteDate.User = user.Name
			tempalteDate.IsAuthenticated = true
		}
	}
	posts, err := c.service.Post.AllPost()
	if err != nil {
		pkg.ResponseError(w, "Status Internal ServerError", http.StatusInternalServerError)
		return
	}
	tempalteDate.Posts = posts
	tmpl, err := template.ParseFiles("./templates/html/home.page.html")
	if err != nil {
		pkg.ResponseError(w, "Status Internal ServerError", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, tempalteDate)
}
