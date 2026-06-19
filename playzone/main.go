package main

import (
	"log"
	"net/http"

	"playzone/backend/config"
	"playzone/backend/database"
	"playzone/backend/router"
)

func main() {
	database.Init()
	defer database.DB.Close()

	r := router.New()

	log.Println("PLAYZONE en route sur http://localhost" + config.Port)
	if err := http.ListenAndServe(config.Port, r); err != nil {
		log.Fatal(err)
	}
}
