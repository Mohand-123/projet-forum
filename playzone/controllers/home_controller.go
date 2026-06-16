package controllers

import (
	"net/http"

	"playzone/middleware"
	"playzone/repositories"
)

// Home affiche la page d'accueil
func Home(w http.ResponseWriter, r *http.Request) {
	categories, _ := repositories.ListAllCategories()
	threads, _ := repositories.ListRecentThreads(10)

	Render(w, r, "home", PageData{
		Title: "Accueil",
		Data: map[string]interface{}{
			"Categories": categories,
			"Threads":    threads,
		},
	})
}

// Search recherche dans les fils (auth requise)
func Search(w http.ResponseWriter, r *http.Request) {
	if middleware.GetUserID(r) == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	q := r.URL.Query().Get("q")
	var threads []interface{}
	if q != "" {
		res, _ := repositories.SearchThreads(q)
		for _, t := range res {
			threads = append(threads, t)
		}
	}

	Render(w, r, "search", PageData{
		Title: "Recherche",
		Data: map[string]interface{}{
			"Query":   q,
			"Threads": threads,
		},
	})
}
