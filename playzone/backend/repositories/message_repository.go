package repositories

import (
	"fmt"

	"playzone/backend/database"
	"playzone/backend/models"
)

// CreateMessage cree un message dans un fil
func CreateMessage(threadID, authorID int, content string) (int64, error) {
	res, err := database.DB.Exec(`
		INSERT INTO messages (thread_id, author_id, content) VALUES (?, ?, ?)
	`, threadID, authorID, content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FindMessageByID retourne un message
func FindMessageByID(id int) (*models.Message, error) {
	m := &models.Message{}
	err := database.DB.QueryRow(`
		SELECT id, thread_id, author_id, content, created_at, updated_at
		FROM messages WHERE id = ?
	`, id).Scan(&m.ID, &m.ThreadID, &m.AuthorID, &m.Content, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// UpdateMessage modifie le contenu
func UpdateMessage(id int, content string) error {
	_, err := database.DB.Exec(`
		UPDATE messages SET content = ?, updated_at = datetime('now') WHERE id = ?
	`, content, id)
	return err
}

// DeleteMessage supprime un message (cascade pour les reactions)
func DeleteMessage(id int) error {
	_, err := database.DB.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}

// ListMessagesByThread liste les messages d'un fil avec tri et pagination
// sortBy = "date" (chronologique) ou "score" (popularite)
// limit = 0 -> pas de limite
func ListMessagesByThread(threadID int, sortBy string, limit, offset int, currentUserID int) ([]models.Message, error) {
	orderBy := "m.created_at DESC"
	if sortBy == "date_asc" {
		orderBy = "m.created_at ASC"
	} else if sortBy == "score" {
		orderBy = "score DESC, m.created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT m.id, m.thread_id, m.author_id, m.content, m.created_at, m.updated_at,
		       u.username,
		       COALESCE(SUM(CASE WHEN r.type = 'like' THEN 1 WHEN r.type = 'dislike' THEN -1 ELSE 0 END), 0) AS score,
		       COALESCE((SELECT type FROM reactions WHERE user_id = ? AND message_id = m.id), '') AS my_react
		FROM messages m
		JOIN users u ON u.id = m.author_id
		LEFT JOIN reactions r ON r.message_id = m.id
		WHERE m.thread_id = ?
		GROUP BY m.id
		ORDER BY %s
	`, orderBy)

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := database.DB.Query(query, currentUserID, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		rows.Scan(&m.ID, &m.ThreadID, &m.AuthorID, &m.Content, &m.CreatedAt, &m.UpdatedAt,
			&m.AuthorName, &m.Score, &m.UserReact)
		messages = append(messages, m)
	}
	return messages, nil
}

// CountMessagesByThread pour la pagination
func CountMessagesByThread(threadID int) int {
	var n int
	database.DB.QueryRow("SELECT COUNT(*) FROM messages WHERE thread_id = ?", threadID).Scan(&n)
	return n
}
