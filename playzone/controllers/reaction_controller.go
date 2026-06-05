package controllers

import (
	"net/http"
	"strconv"

	"playzone/middleware"
	"playzone/repositories"
	"playzone/services"
)

// React applique une reaction
func React(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	reactionType := r.FormValue("type")

	err := services.React(middleware.GetUserID(r), messageID, reactionType)
	if err != nil {
		SetFlash(w, "error", err.Error())
	}

	// Redirection vers le fil concerne
	m, _ := repositories.FindMessageByID(messageID)
	if m != nil {
		http.Redirect(w, r, "/threads/"+strconv.Itoa(m.ThreadID)+"#msg-"+strconv.Itoa(messageID), http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
