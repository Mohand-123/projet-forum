package config

// Configuration globale de l'application
const (
	Port      = ":3000"
	DBPath    = "data/playzone.db"
	JWTSecret = "playzone-secret-key-a-changer-en-prod"
	JWTHours  = 24 * 7 // duree du token en heures (7 jours)
)
