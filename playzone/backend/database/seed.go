package database

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"log"
)

// hashPassword genere un salt et hash le mdp en SHA-512
func hashPassword(password string) (string, string) {
	saltBytes := make([]byte, 16)
	rand.Read(saltBytes)
	salt := hex.EncodeToString(saltBytes)

	h := sha512.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil)), salt
}

// Seed insere les donnees de base si la table users est vide
func Seed() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count > 0 {
		return
	}

	log.Println("Seed de la base avec donnees de demo...")

	// Comptes par defaut
	adminHash, adminSalt := hashPassword("Admin1234567!")
	demoHash, demoSalt := hashPassword("Demo1234567!")

	DB.Exec(`INSERT INTO users (username, email, password_hash, salt, role) VALUES (?, ?, ?, ?, ?)`,
		"admin", "admin@playzone.dz", adminHash, adminSalt, "admin")
	DB.Exec(`INSERT INTO users (username, email, password_hash, salt, role) VALUES (?, ?, ?, ?, ?)`,
		"demo", "demo@playzone.dz", demoHash, demoSalt, "user")

	// Categories
	categories := []struct {
		name, slug string
	}{
		{"Counter-Strike 2", "cs2"},
		{"League of Legends", "lol"},
		{"Valorant", "valorant"},
		{"Fortnite", "fortnite"},
		{"Minecraft", "minecraft"},
		{"General", "general"},
	}
	for _, c := range categories {
		DB.Exec("INSERT INTO categories (name, slug) VALUES (?, ?)", c.name, c.slug)
	}

	// Tags
	tags := []struct {
		name, color string
	}{
		{"strategie", "#7c3aed"},
		{"clip", "#06b6d4"},
		{"guide", "#10b981"},
		{"question", "#f59e0b"},
		{"annonce", "#ef4444"},
		{"lfg", "#ec4899"},
	}
	for _, t := range tags {
		DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", t.name, t.color)
	}

	// Fil d'exemple
	res, _ := DB.Exec(`INSERT INTO threads (title, content, author_id, category_id, state)
		VALUES (?, ?, ?, ?, ?)`,
		"Bienvenue sur PLAYZONE",
		"Hello la communaute ! Ce fil est la pour souhaiter la bienvenue a tous les nouveaux membres. Presentez-vous, dites quels jeux vous aimez. GG !",
		1, 6, "ouvert")
	threadID, _ := res.LastInsertId()

	DB.Exec("INSERT INTO thread_tags (thread_id, tag_id) VALUES (?, ?)", threadID, 5) // annonce

	DB.Exec(`INSERT INTO messages (thread_id, author_id, content) VALUES (?, ?, ?)`,
		threadID, 2, "Yo ! Ravi d'etre la, j'attends de discuter strategies !")

	log.Println("Seed termine. Comptes : admin/Admin1234567! et demo/Demo1234567!")
}
