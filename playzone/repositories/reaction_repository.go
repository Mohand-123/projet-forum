package repositories

import "playzone/database"

// FindReaction cherche la reaction d'un user sur un message
func FindReaction(userID, messageID int) (string, bool) {
	var t string
	err := database.DB.QueryRow(`
		SELECT type FROM reactions WHERE user_id = ? AND message_id = ?
	`, userID, messageID).Scan(&t)
	if err != nil {
		return "", false
	}
	return t, true
}

// AddReaction insere une nouvelle reaction
func AddReaction(userID, messageID int, reactionType string) error {
	_, err := database.DB.Exec(`
		INSERT INTO reactions (user_id, message_id, type) VALUES (?, ?, ?)
	`, userID, messageID, reactionType)
	return err
}

// UpdateReaction change le type (like <-> dislike)
func UpdateReaction(userID, messageID int, reactionType string) error {
	_, err := database.DB.Exec(`
		UPDATE reactions SET type = ?, created_at = datetime('now')
		WHERE user_id = ? AND message_id = ?
	`, reactionType, userID, messageID)
	return err
}

// RemoveReaction supprime la reaction
func RemoveReaction(userID, messageID int) error {
	_, err := database.DB.Exec(`
		DELETE FROM reactions WHERE user_id = ? AND message_id = ?
	`, userID, messageID)
	return err
}
