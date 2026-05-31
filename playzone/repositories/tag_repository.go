package repositories

import (
	"playzone/database"
	"playzone/models"
)

// ListAllCategories retourne toutes les categories
func ListAllCategories() ([]models.Category, error) {
	rows, err := database.DB.Query("SELECT id, name, slug FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name, &c.Slug)
		cats = append(cats, c)
	}
	return cats, nil
}

// FindCategoryBySlug retourne une categorie via son slug
func FindCategoryBySlug(slug string) (*models.Category, error) {
	c := &models.Category{}
	err := database.DB.QueryRow(`
		SELECT id, name, slug FROM categories WHERE slug = ?
	`, slug).Scan(&c.ID, &c.Name, &c.Slug)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ListAllTags retourne tous les tags
func ListAllTags() ([]models.Tag, error) {
	rows, err := database.DB.Query("SELECT id, name, color FROM tags ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		rows.Scan(&t.ID, &t.Name, &t.Color)
		tags = append(tags, t)
	}
	return tags, nil
}
