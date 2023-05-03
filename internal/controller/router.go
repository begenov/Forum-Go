package controller

import "net/http"

func (c *Controller) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", c.middleware(c.home))

	mux.HandleFunc("/auth/sign-up", c.middleware(c.User.SignUp))
	mux.HandleFunc("/auth/sign-in", c.middleware(c.User.SignIn))
	mux.HandleFunc("/auth/log-out", c.User.LogOut)

	mux.HandleFunc("/post/create", c.middleware(c.Post.CreatePost))
	mux.HandleFunc("/post", c.middleware(c.Post.Post))
	mux.HandleFunc("/post/mylikedpost", c.middleware(c.Post.MyLikedPost))
	mux.HandleFunc("/post/mypost", c.middleware(c.Post.MyPost))

	mux.HandleFunc("/post/category", c.middleware(c.Post.PostCategory))

	mux.HandleFunc("/reaction/post", c.middleware(c.Reaction.ReactionPost))
	mux.HandleFunc("/reaction/comment", c.middleware(c.Reaction.ReactionComment))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	return mux
}
