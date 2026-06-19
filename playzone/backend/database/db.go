package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"playzone/backend/config"
)

var DB *sql.DB

// Init ouvre la connexion et applique le schema
func Init() {
	// On verifie que le dossier data existe
	os.MkdirAll("data", 0755)

	var err error
	DB, err = sql.Open("sqlite", config.DBPath)
	if err != nil {
		log.Fatal("Impossible d'ouvrir la BDD :", err)
	}

	// Active les foreign keys
	DB.Exec("PRAGMA foreign_keys = ON;")

	// Applique le schema
	schema, err := os.ReadFile("backend/database/schema.sql")
	if err != nil {
		log.Fatal("Schema introuvable :", err)
	}
	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Erreur creation tables :", err)
	}

	// Seed initial si la base est vide
	Seed()
}
