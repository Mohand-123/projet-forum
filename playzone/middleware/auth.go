package middleware

import (
	"context"
	"net/http"

	"playzone/services"
)

type ctxKey string

const (
	CtxUserID   ctxKey = "user_id"
	CtxUsername ctxKey = "username"
	CtxRole     ctxKey = "role"
)

// getUserFromCookie extrait les infos du JWT stocke en cookie
func getUserFromCookie(r *http.Request) (int, string, string, bool) {
	cookie, err := r.Cookie("playzone_token")
	if err != nil {
		return 0, "", "", false
	}
	claims, err := services.ParseJWT(cookie.Value)
	if err != nil {
		return 0, "", "", false
	}
	uid := int(claims["user_id"].(float64))
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)
	return uid, username, role, true
}

// AttachUser ajoute les infos user au contexte si connecte (pas de blocage)
func AttachUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, username, role, ok := getUserFromCookie(r)
		if ok {
			ctx := context.WithValue(r.Context(), CtxUserID, uid)
			ctx = context.WithValue(ctx, CtxUsername, username)
			ctx = context.WithValue(ctx, CtxRole, role)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAuth bloque si pas connecte
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(CtxUserID).(int); !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin bloque si pas admin
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(CtxRole).(string)
		if !ok || role != "admin" {
			http.Error(w, "Acces refuse", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserID helper
func GetUserID(r *http.Request) int {
	if v, ok := r.Context().Value(CtxUserID).(int); ok {
		return v
	}
	return 0
}

// GetUsername helper
func GetUsername(r *http.Request) string {
	if v, ok := r.Context().Value(CtxUsername).(string); ok {
		return v
	}
	return ""
}

// IsAdmin helper
func IsAdmin(r *http.Request) bool {
	if v, ok := r.Context().Value(CtxRole).(string); ok {
		return v == "admin"
	}
	return false
}
