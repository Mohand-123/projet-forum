package repositories

import (
	"playzone/database"
	"playzone/models"
)

// FindUserByLogin cherche un user par username OU email
func FindUserByLogin(login string) (*models.User, error) {
	u := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, salt, role, is_banned, created_at
		FROM users WHERE username = ? OR email = ?
	`, login, login).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Salt, &u.Role, &u.IsBanned, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindUserByID cherche un user par son id
func FindUserByID(id int) (*models.User, error) {
	u := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, salt, role, is_banned, created_at
		FROM users WHERE id = ?
	`, id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Salt, &u.Role, &u.IsBanned, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// CreateUser cree un nouveau compte
func CreateUser(username, email, passwordHash, salt string) (int64, error) {
	res, err := database.DB.Exec(`
		INSERT INTO users (username, email, password_hash, salt) VALUES (?, ?, ?, ?)
	`, username, email, passwordHash, salt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ExistsByUsernameOrEmail verifie si username ou email est deja pris
func ExistsByUsernameOrEmail(username, email string) bool {
	var count int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE username = ? OR email = ?
	`, username, email).Scan(&count)
	return count > 0
}

// ListAllUsers retourne tous les users (pour admin)
func ListAllUsers() ([]models.User, error) {
	rows, err := database.DB.Query(`
		SELECT id, username, email, role, is_banned, created_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.IsBanned, &u.CreatedAt)
		users = append(users, u)
	}
	return users, nil
}

// SetBanned modifie le statut banni d'un user
func SetBanned(userID int, banned bool) error {
	val := 0
	if banned {
		val = 1
	}
	_, err := database.DB.Exec("UPDATE users SET is_banned = ? WHERE id = ?", val, userID)
	return err
}
