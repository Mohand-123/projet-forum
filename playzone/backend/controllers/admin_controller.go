package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"playzone/backend/repositories"
	"playzone/backend/services"
)

// AdminDashboard affiche les stats globales
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	users, _ := repositories.ListAllUsers()

	var nbThreads, nbMessages int
	repositories.GetDB().QueryRow("SELECT COUNT(*) FROM threads").Scan(&nbThreads)
	repositories.GetDB().QueryRow("SELECT COUNT(*) FROM messages").Scan(&nbMessages)

	Render(w, r, "admin", PageData{
		Title: "Administration",
		Data: map[string]interface{}{
			"Users":      users,
			"NbThreads":  nbThreads,
			"NbMessages": nbMessages,
		},
	})
}

// BanUser bannit un user
func BanUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := services.BanUser(id); err != nil {
		SetFlash(w, "error", err.Error())
	} else {
		SetFlash(w, "success", "Utilisateur banni")
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// UnbanUser deban
func UnbanUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := services.UnbanUser(id); err != nil {
		SetFlash(w, "error", err.Error())
	} else {
		SetFlash(w, "success", "Utilisateur deban")
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}
