package services

import (
	"errors"

	"playzone/backend/repositories"
)

// CreateMessage poste un message dans un fil (verif que le fil est ouvert)
func CreateMessage(threadID, authorID int, content string) (int64, error) {
	if len(content) < 1 {
		return 0, errors.New("le message ne peut pas etre vide")
	}

	t, err := repositories.FindThreadByID(threadID)
	if err != nil {
		return 0, errors.New("fil introuvable")
	}
	if t.State != "ouvert" {
		return 0, errors.New("ce fil n'accepte plus de nouveaux messages")
	}

	return repositories.CreateMessage(threadID, authorID, content)
}

// UpdateMessage modifie un message (verif auteur)
func UpdateMessage(messageID, userID int, isAdmin bool, content string) error {
	m, err := repositories.FindMessageByID(messageID)
	if err != nil {
		return errors.New("message introuvable")
	}
	if m.AuthorID != userID && !isAdmin {
		return errors.New("tu n'as pas le droit de modifier ce message")
	}
	if len(content) < 1 {
		return errors.New("contenu vide")
	}
	return repositories.UpdateMessage(messageID, content)
}

// DeleteMessage supprime un message (auteur ou admin)
func DeleteMessage(messageID, userID int, isAdmin bool) error {
	m, err := repositories.FindMessageByID(messageID)
	if err != nil {
		return errors.New("message introuvable")
	}
	if m.AuthorID != userID && !isAdmin {
		return errors.New("tu n'as pas le droit de supprimer ce message")
	}
	return repositories.DeleteMessage(messageID)
}
