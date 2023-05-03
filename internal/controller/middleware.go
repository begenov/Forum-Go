package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/begenov/Forum-Go/models"
)

func (c *Controller) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx := context.WithValue(r.Context(), models.UserIDKey, 0)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if err == cookie.Valid() {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		flag := false
		user, err := c.service.User.UserByToken(cookie.Value)
		if err != nil {
			flag = true
			ctx := context.WithValue(r.Context(), models.UserIDKey, 0)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if user.ExpiresAt.Before(time.Now()) {
			flag = true
			if err := c.service.User.DeleteSession(user); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/user/sign-in", http.StatusSeeOther)
			return
		}
		if flag {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "",
				MaxAge:   -1,
				HttpOnly: true,
			})
		}
		ctx := context.WithValue(r.Context(), models.UserIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// func (c *Controller) auth(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, err := r.Cookie("session")

// 		if err == http.ErrNoCookie {
// 			ctx := context.WithValue(r.Context(), models.UserIDKey, true)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 			return
// 		}

// 		http.Error(w, "выйди из аккаунта", http.StatusBadRequest)
// 	}
// }
