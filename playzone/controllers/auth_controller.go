package controllers

import (
	"net/http"

	"playzone/config"
	"playzone/repositories"
	"playzone/services"
)

// GetRegister affiche le formulaire d'inscription
func GetRegister(w http.ResponseWriter, r *http.Request) {
	Render(w, r, "register", PageData{Title: "Inscription"})
}

// PostRegister traite le formulaire
func PostRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	uid, err := services.Register(username, email, password)
	if err != nil {
		Render(w, r, "register", PageData{
			Title: "Inscription",
			Error: err.Error(),
			Data:  map[string]interface{}{"Username": username, "Email": email},
		})
		return
	}

	// Auto-connexion
	user, _ := repositories.FindUserByID(int(uid))
	token, _ := services.GenerateJWT(user)
	http.SetCookie(w, &http.Cookie{
		Name:     "playzone_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   config.JWTHours * 3600,
	})
	SetFlash(w, "success", "Bienvenue sur PLAYZONE !")
	http.Redirect(w, r, "/", http.StatusFound)
}

// GetLogin affiche le formulaire
func GetLogin(w http.ResponseWriter, r *http.Request) {
	Render(w, r, "login", PageData{Title: "Connexion"})
}

// PostLogin traite la connexion
func PostLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	login := r.FormValue("login")
	password := r.FormValue("password")

	user, err := services.Login(login, password)
	if err != nil {
		Render(w, r, "login", PageData{
			Title: "Connexion",
			Error: err.Error(),
			Data:  map[string]interface{}{"Login": login},
		})
		return
	}

	token, err := services.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Erreur generation token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "playzone_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   config.JWTHours * 3600,
	})
	SetFlash(w, "success", "Connecte en tant que "+user.Username)
	http.Redirect(w, r, "/", http.StatusFound)
}

// PostLogout supprime le cookie
func PostLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "playzone_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
