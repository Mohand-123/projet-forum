package repositories

import (
	"fmt"
	"strings"

	"playzone/backend/database"
	"playzone/backend/models"
)

// CreateThread cree un nouveau fil
func CreateThread(title, content string, authorID, categoryID int) (int64, error) {
	res, err := database.DB.Exec(`
		INSERT INTO threads (title, content, author_id, category_id) VALUES (?, ?, ?, ?)
	`, title, content, authorID, categoryID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// AttachTags associe des tags a un fil
func AttachTags(threadID int64, tagIDs []int) error {
	for _, id := range tagIDs {
		_, err := database.DB.Exec("INSERT OR IGNORE INTO thread_tags (thread_id, tag_id) VALUES (?, ?)", threadID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// FindThreadByID retourne un fil avec ses infos
func FindThreadByID(id int) (*models.Thread, error) {
	t := &models.Thread{}
	err := database.DB.QueryRow(`
		SELECT t.id, t.title, t.content, t.author_id, t.category_id, t.state, t.created_at, t.updated_at,
		       u.username, c.name, c.slug
		FROM threads t
		JOIN users u      ON u.id = t.author_id
		JOIN categories c ON c.id = t.category_id
		WHERE t.id = ?
	`, id).Scan(&t.ID, &t.Title, &t.Content, &t.AuthorID, &t.CategoryID, &t.State,
		&t.CreatedAt, &t.UpdatedAt, &t.AuthorName, &t.CategoryName, &t.CategorySlug)
	if err != nil {
		return nil, err
	}

	// Tags
	rows, _ := database.DB.Query(`
		SELECT t.id, t.name, t.color FROM tags t
		JOIN thread_tags tt ON tt.tag_id = t.id
		WHERE tt.thread_id = ?
	`, id)
	defer rows.Close()
	for rows.Next() {
		var tag models.Tag
		rows.Scan(&tag.ID, &tag.Name, &tag.Color)
		t.Tags = append(t.Tags, tag)
	}

	return t, nil
}

// UpdateThread modifie titre et contenu (pas l'etat ni les tags ici)
func UpdateThread(id int, title, content string) error {
	_, err := database.DB.Exec(`
		UPDATE threads SET title = ?, content = ?, updated_at = datetime('now')
		WHERE id = ?
	`, title, content, id)
	return err
}

// UpdateThreadState change l'etat (ouvert / ferme / archive)
func UpdateThreadState(id int, state string) error {
	_, err := database.DB.Exec("UPDATE threads SET state = ? WHERE id = ?", state, id)
	return err
}

// DeleteThread supprime un fil (et cascade messages, reactions, thread_tags)
func DeleteThread(id int) error {
	_, err := database.DB.Exec("DELETE FROM threads WHERE id = ?", id)
	return err
}

// ListThreadsByCategory liste les fils d'une categorie (state != archive)
// limit = 0 signifie pas de limite (toutes les donnees)
func ListThreadsByCategory(categorySlug string, limit, offset int) ([]models.Thread, error) {
	query := `
		SELECT t.id, t.title, t.state, t.created_at, t.author_id,
		       u.username, c.name, c.slug,
		       (SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id) AS nb_msg
		FROM threads t
		JOIN users u      ON u.id = t.author_id
		JOIN categories c ON c.id = t.category_id
		WHERE c.slug = ? AND t.state != 'archive'
		ORDER BY t.created_at DESC
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := database.DB.Query(query, categorySlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		rows.Scan(&t.ID, &t.Title, &t.State, &t.CreatedAt, &t.AuthorID,
			&t.AuthorName, &t.CategoryName, &t.CategorySlug, &t.MessageCount)
		threads = append(threads, t)
	}
	return threads, nil
}

// CountThreadsByCategory pour la pagination
func CountThreadsByCategory(slug string) int {
	var n int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM threads t
		JOIN categories c ON c.id = t.category_id
		WHERE c.slug = ? AND t.state != 'archive'
	`, slug).Scan(&n)
	return n
}

// ListRecentThreads pour la page d'accueil
func ListRecentThreads(limit int) ([]models.Thread, error) {
	rows, err := database.DB.Query(`
		SELECT t.id, t.title, t.state, t.created_at, t.author_id,
		       u.username, c.name, c.slug,
		       (SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id) AS nb_msg
		FROM threads t
		JOIN users u      ON u.id = t.author_id
		JOIN categories c ON c.id = t.category_id
		WHERE t.state != 'archive'
		ORDER BY t.created_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		rows.Scan(&t.ID, &t.Title, &t.State, &t.CreatedAt, &t.AuthorID,
			&t.AuthorName, &t.CategoryName, &t.CategorySlug, &t.MessageCount)
		threads = append(threads, t)
	}
	return threads, nil
}

// SearchThreads cherche par titre ou par tag/categorie
func SearchThreads(q string) ([]models.Thread, error) {
	like := "%" + q + "%"
	rows, err := database.DB.Query(`
		SELECT DISTINCT t.id, t.title, t.state, t.created_at,
		       u.username, c.name, c.slug
		FROM threads t
		JOIN users u      ON u.id = t.author_id
		JOIN categories c ON c.id = t.category_id
		LEFT JOIN thread_tags tt ON tt.thread_id = t.id
		LEFT JOIN tags tag       ON tag.id = tt.tag_id
		WHERE t.state != 'archive'
		  AND (t.title LIKE ? OR c.name LIKE ? OR c.slug LIKE ? OR tag.name LIKE ?)
		ORDER BY t.created_at DESC LIMIT 50
	`, like, like, like, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		rows.Scan(&t.ID, &t.Title, &t.State, &t.CreatedAt,
			&t.AuthorName, &t.CategoryName, &t.CategorySlug)
		threads = append(threads, t)
	}
	return threads, nil
}

// ListThreadsByTag retourne les fils ayant un tag particulier
func ListThreadsByTag(tagName string) ([]models.Thread, error) {
	rows, err := database.DB.Query(`
		SELECT t.id, t.title, t.state, t.created_at,
		       u.username, c.name, c.slug
		FROM threads t
		JOIN users u       ON u.id = t.author_id
		JOIN categories c  ON c.id = t.category_id
		JOIN thread_tags tt ON tt.thread_id = t.id
		JOIN tags tag       ON tag.id = tt.tag_id
		WHERE tag.name = ? AND t.state != 'archive'
		ORDER BY t.created_at DESC
	`, tagName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		rows.Scan(&t.ID, &t.Title, &t.State, &t.CreatedAt,
			&t.AuthorName, &t.CategoryName, &t.CategorySlug)
		threads = append(threads, t)
	}
	return threads, nil
}

// ReplaceTags remplace tous les tags d'un fil par la nouvelle liste
func ReplaceTags(threadID int, tagIDs []int) error {
	database.DB.Exec("DELETE FROM thread_tags WHERE thread_id = ?", threadID)
	for _, id := range tagIDs {
		database.DB.Exec("INSERT OR IGNORE INTO thread_tags (thread_id, tag_id) VALUES (?, ?)", threadID, id)
	}
	return nil
}

// helper pour eviter import inutile
var _ = strings.HasPrefix
