package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"playzone/backend/controllers"
	"playzone/backend/middleware"
)

// New cree le routeur principal
func New() http.Handler {
	r := chi.NewRouter()

	// Middleware global : attache l'user au contexte si JWT present
	r.Use(middleware.AttachUser)

	// Fichiers statiques (CSS, JS)
	fs := http.FileServer(http.Dir("frontend/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Routes publiques
	r.Get("/", controllers.Home)
	r.Get("/register", controllers.GetRegister)
	r.Post("/register", controllers.PostRegister)
	r.Get("/login", controllers.GetLogin)
	r.Post("/login", controllers.PostLogin)
	r.Post("/logout", controllers.PostLogout)

	r.Get("/categories/{slug}", controllers.ShowCategory)
	r.Get("/tags/{name}", controllers.ShowByTag)
	r.Get("/threads/{id:[0-9]+}", controllers.ShowThread)

	// Recherche : auth requise
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/search", controllers.Search)

		// Threads
		r.Get("/threads/new", controllers.NewThreadForm)
		r.Post("/threads", controllers.CreateThread)
		r.Get("/threads/{id}/edit", controllers.EditThreadForm)
		r.Post("/threads/{id}/update", controllers.UpdateThread)
		r.Post("/threads/{id}/delete", controllers.DeleteThread)
		r.Post("/threads/{id}/state", controllers.ChangeState)

		// Messages
		r.Post("/messages", controllers.CreateMessage)
		r.Post("/messages/{id}/update", controllers.UpdateMessage)
		r.Post("/messages/{id}/delete", controllers.DeleteMessage)

		// Reactions
		r.Post("/reactions", controllers.React)
	})

	// Admin
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Use(middleware.RequireAdmin)
		r.Get("/admin", controllers.AdminDashboard)
		r.Post("/admin/users/{id}/ban", controllers.BanUser)
		r.Post("/admin/users/{id}/unban", controllers.UnbanUser)
	})

	return r
}
