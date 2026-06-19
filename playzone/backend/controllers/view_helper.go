package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"playzone/backend/middleware"
)

// PageData est la structure passee a toutes les vues
type PageData struct {
	Title       string
	UserID      int
	Username    string
	IsAdmin     bool
	IsLogged    bool
	Flash       string
	FlashType   string
	Error       string
	CurrentPath string
	Data        map[string]interface{}
}

// Render charge le layout + la page et l'execute
func Render(w http.ResponseWriter, r *http.Request, page string, data PageData) {
	// Auto-remplit les infos user
	data.UserID = middleware.GetUserID(r)
	data.Username = middleware.GetUsername(r)
	data.IsAdmin = middleware.IsAdmin(r)
	data.IsLogged = data.UserID > 0
	data.CurrentPath = r.URL.Path

	// Recupere le flash depuis le cookie
	if cookie, err := r.Cookie("flash"); err == nil {
		data.Flash = cookie.Value
		if t, err := r.Cookie("flash_type"); err == nil {
			data.FlashType = t.Value
		}
		// On efface le flash
		http.SetCookie(w, &http.Cookie{Name: "flash", Path: "/", MaxAge: -1})
		http.SetCookie(w, &http.Cookie{Name: "flash_type", Path: "/", MaxAge: -1})
	}

	layout := filepath.Join("frontend", "views", "layout.html")
	pagePath := filepath.Join("frontend", "views", page+".html")

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"seq": func(start, end int) []int {
			r := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				r = append(r, i)
			}
			return r
		},
	}

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles(layout, pagePath)
	if err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur execution : "+err.Error(), http.StatusInternalServerError)
	}
}

// SetFlash stocke un message flash dans un cookie
func SetFlash(w http.ResponseWriter, flashType, message string) {
	http.SetCookie(w, &http.Cookie{Name: "flash", Value: message, Path: "/", MaxAge: 30})
	http.SetCookie(w, &http.Cookie{Name: "flash_type", Value: flashType, Path: "/", MaxAge: 30})
}
