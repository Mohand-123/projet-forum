package services

import (
	"errors"

	"playzone/backend/repositories"
)

// React applique un like ou dislike (toggle si meme type)
func React(userID, messageID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return errors.New("reaction invalide")
	}

	existing, found := repositories.FindReaction(userID, messageID)
	if !found {
		return repositories.AddReaction(userID, messageID, reactionType)
	}

	if existing == reactionType {
		// Meme reaction -> on annule
		return repositories.RemoveReaction(userID, messageID)
	}

	// Reaction differente -> on update
	return repositories.UpdateReaction(userID, messageID, reactionType)
}
