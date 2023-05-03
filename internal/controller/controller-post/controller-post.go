package controllerpost

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/begenov/Forum-Go/models"
	"github.com/begenov/Forum-Go/pkg"
)

type postProvider interface {
	CreatePost(post models.Post) (int, error)
	GetPostByID(id int) (models.Post, error)
	AllPost() ([]models.Post, error)
	GetPostsByUserID(userID int) ([]models.Post, error)
	GetMyLikedPosts(userID int) ([]models.Post, error)
	GetPostsByCategory(categroy []string) ([]models.Post, error)
}

type userProvider interface {
	UserByID(id int) (models.User, error)
}

type commentProvider interface {
	CreateComment(c models.Comment) error
	GetCommentByPostId(id int) ([]models.Comment, error)
}

type PostController struct {
	postProvider postProvider
	userProvider userProvider
	comment      commentProvider
}

func NewPostController(postProvider postProvider, userpostProvider userProvider, commentProvider commentProvider) *PostController {
	return &PostController{
		postProvider: postProvider, userProvider: userpostProvider,
		comment: commentProvider,
	}
}

func (c *PostController) Post(w http.ResponseWriter, r *http.Request) {
	var tempalteDate models.TemplateData
	userID := r.Context().Value(models.UserIDKey)

	if userID.(int) <= 0 {
		tempalteDate.IsAuthenticated = false
	}

	user, _ := c.userProvider.UserByID(userID.(int))
	if user.ID == 0 {
		tempalteDate.IsAuthenticated = false
	} else {
		tempalteDate.IsAuthenticated = true
		tempalteDate.User = user.Name
	}
	postID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {

		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	post, err := c.postProvider.GetPostByID(postID)
	if err != nil || post.Id == 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	comments, err := c.comment.GetCommentByPostId(postID)
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tempalteDate.Comment = comments
	tempalteDate.Post = post
	switch r.Method {
	case http.MethodPost:
		comment := getcomment(r, tempalteDate, post)

		if err := c.comment.CreateComment(comment); err != nil {
			tempalteDate.MsgError = "Empty text error"
			tmpl, err := template.ParseFiles("./templates/html/post.page.html")
			if err != nil {
				pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, tempalteDate)
			return
		}
		http.Redirect(w, r, r.URL.Path+"?id="+r.FormValue("id"), http.StatusSeeOther)
	case http.MethodGet:
		tmpl, err := template.ParseFiles("./templates/html/post.page.html")
		if err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, tempalteDate)
	default:
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var tempalteDate models.TemplateData
	userID := r.Context().Value(models.UserIDKey)

	if userID.(int) <= 0 {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	user, err := c.userProvider.UserByID(userID.(int))
	if err != nil || user.ID == 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	tempalteDate.User = user.Name
	tempalteDate.IsAuthenticated = true
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("./templates/html/create-post.html")
		if err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, tempalteDate)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		post := parsepost(r, user)
		postID, err := c.postProvider.CreatePost(post)
		if err != nil {
			tempalteDate.MsgError = err.Error()
			tmpl, err := template.ParseFiles("./templates/html/create-post.html")
			if err != nil {
				pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			tmpl.Execute(w, tempalteDate)

		}
		page := fmt.Sprintf("/post?id=%d", postID)
		http.Redirect(w, r, page, http.StatusSeeOther)
	default:
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (c *PostController) MyLikedPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var tempalteDate models.TemplateData

	userID := r.Context().Value(models.UserIDKey)
	if userID.(int) <= 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := c.userProvider.UserByID(userID.(int))
	if err != nil || user.ID == 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tempalteDate.User = user.Name
	posts, err := c.postProvider.GetMyLikedPosts(user.ID)
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tempalteDate.IsAuthenticated = true
	tempalteDate.Posts = posts

	tmpl, err := template.ParseFiles("./templates/html/posts.html")
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, tempalteDate)
}

func (c *PostController) MyPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value(models.UserIDKey)
	if userID.(int) <= 0 {
		pkg.ResponseError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := c.userProvider.UserByID(userID.(int))
	if err != nil || user.ID == 0 {
		pkg.ResponseError(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	var tempalteDate models.TemplateData

	tempalteDate.IsAuthenticated = true
	tempalteDate.User = user.Name

	posts, err := c.postProvider.GetPostsByUserID(user.ID)
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tempalteDate.Posts = posts
	tmpl, err := template.ParseFiles("./templates/html/posts.html")
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, tempalteDate)
}

func (c *PostController) PostCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.ResponseError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var tempalteDate models.TemplateData
	if err := r.ParseForm(); err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userID := r.Context().Value(models.UserIDKey)

	if userID.(int) <= 0 {
		tempalteDate.IsAuthenticated = false
	} else {
		user, _ := c.userProvider.UserByID(userID.(int))
		if user.ID == 0 {
			tempalteDate.IsAuthenticated = false
		} else {
			tempalteDate.IsAuthenticated = true
			tempalteDate.User = user.Name
		}
	}

	values := r.URL.Query()
	categories := values["category"]
	posts, err := c.postProvider.GetPostsByCategory(categories)
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	tempalteDate.Posts = posts
	tmpl, err := template.ParseFiles("./templates/html/categories.page.html")
	if err != nil {
		pkg.ResponseError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, tempalteDate)
}

func parsepost(r *http.Request, u models.User) models.Post {
	title := r.FormValue("title")
	category := r.Form["category"]
	content := r.FormValue("content")
	return models.Post{
		Title:       title,
		Category:    category,
		Description: content,
		CreateAt:    time.Now().Format("January 2, 2006"),
		Author:      u.Name,
		UserID:      u.ID,
	}
}

func getcomment(r *http.Request, tempalteDate models.TemplateData, post models.Post) models.Comment {
	text := r.FormValue("comment-text")
	return models.Comment{
		UserID: tempalteDate.UserID,
		PostID: post.Id,
		Author: tempalteDate.User,
		Text:   text,
		Date:   time.Now().Format("01-02-2006 15:04:05"),
	}
}
