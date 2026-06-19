package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"playzone/backend/middleware"
	"playzone/backend/repositories"
	"playzone/backend/services"
)

// CreateMessage poste un nouveau message
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	threadID, _ := strconv.Atoi(r.FormValue("thread_id"))
	content := r.FormValue("content")

	_, err := services.CreateMessage(threadID, middleware.GetUserID(r), content)
	if err != nil {
		SetFlash(w, "error", err.Error())
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
}

// UpdateMessage modifie un message
func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	r.ParseForm()
	content := r.FormValue("content")

	err := services.UpdateMessage(id, middleware.GetUserID(r), middleware.IsAdmin(r), content)
	if err != nil {
		SetFlash(w, "error", err.Error())
	}
	// On retourne dans le fil
	m, _ := repositories.FindMessageByID(id)
	if m != nil {
		http.Redirect(w, r, "/threads/"+strconv.Itoa(m.ThreadID), http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// DeleteMessage supprime
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	m, _ := repositories.FindMessageByID(id)
	threadID := 0
	if m != nil {
		threadID = m.ThreadID
	}

	err := services.DeleteMessage(id, middleware.GetUserID(r), middleware.IsAdmin(r))
	if err != nil {
		SetFlash(w, "error", err.Error())
	} else {
		SetFlash(w, "success", "Message supprime")
	}
	if threadID > 0 {
		http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
