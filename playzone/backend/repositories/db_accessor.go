package repositories

import (
	"database/sql"

	"playzone/backend/database"
)

// GetDB expose la connexion BDD pour les requetes simples des controllers
func GetDB() *sql.DB {
	return database.DB
}
