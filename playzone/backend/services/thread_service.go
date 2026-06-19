package services

import (
	"errors"

	"playzone/backend/models"
	"playzone/backend/repositories"
)

// CreateThread cree un fil + attache les tags
func CreateThread(title, content string, authorID, categoryID int, tagIDs []int) (int64, error) {
	if len(title) < 3 || len(title) > 150 {
		return 0, errors.New("le titre doit faire entre 3 et 150 caracteres")
	}
	if len(content) < 1 {
		return 0, errors.New("le contenu ne peut pas etre vide")
	}

	id, err := repositories.CreateThread(title, content, authorID, categoryID)
	if err != nil {
		return 0, err
	}
	repositories.AttachTags(id, tagIDs)
	return id, nil
}

// GetThread retourne un fil avec verif d'accessibilite (archive bloque)
func GetThread(id int) (*models.Thread, error) {
	t, err := repositories.FindThreadByID(id)
	if err != nil {
		return nil, errors.New("fil introuvable")
	}
	if t.State == "archive" {
		return nil, errors.New("ce fil n'est plus accessible")
	}
	return t, nil
}

// UpdateThread modifie un fil (verif proprietaire)
func UpdateThread(threadID, userID int, isAdmin bool, title, content string, tagIDs []int) error {
	t, err := repositories.FindThreadByID(threadID)
	if err != nil {
		return errors.New("fil introuvable")
	}
	if t.AuthorID != userID && !isAdmin {
		return errors.New("tu n'as pas le droit de modifier ce fil")
	}
	if len(title) < 3 || len(title) > 150 {
		return errors.New("titre invalide")
	}
	if len(content) < 1 {
		return errors.New("contenu vide")
	}
	if err := repositories.UpdateThread(threadID, title, content); err != nil {
		return err
	}
	repositories.ReplaceTags(threadID, tagIDs)
	return nil
}

// DeleteThread supprime un fil (verif proprietaire ou admin)
func DeleteThread(threadID, userID int, isAdmin bool) error {
	t, err := repositories.FindThreadByID(threadID)
	if err != nil {
		return errors.New("fil introuvable")
	}
	if t.AuthorID != userID && !isAdmin {
		return errors.New("tu n'as pas le droit de supprimer ce fil")
	}
	return repositories.DeleteThread(threadID)
}

// ChangeState change l'etat d'un fil (admin only ou proprietaire pour son fil)
func ChangeState(threadID, userID int, isAdmin bool, newState string) error {
	if newState != "ouvert" && newState != "ferme" && newState != "archive" {
		return errors.New("etat invalide")
	}
	t, err := repositories.FindThreadByID(threadID)
	if err != nil {
		return errors.New("fil introuvable")
	}
	if t.AuthorID != userID && !isAdmin {
		return errors.New("non autorise")
	}
	return repositories.UpdateThreadState(threadID, newState)
}
