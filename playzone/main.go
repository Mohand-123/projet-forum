package main

import (
	"log"
	"net/http"

	"playzone/backend/config"
	"playzone/backend/database"
	"playzone/backend/router"
)

// main est le point d'entree de l'application :
// initialise la base de donnees, construit le routeur et demarre le serveur HTTP.
func main() {
	// Ouvre la connexion BDD, applique le schema et le seed
	database.Init()
	defer database.DB.Close()

	// Construit le routeur avec toutes les routes et middlewares
	r := router.New()

	log.Println("PLAYZONE en route sur http://localhost" + config.Port)
	if err := http.ListenAndServe(config.Port, r); err != nil {
		log.Fatal(err)
	}
}
