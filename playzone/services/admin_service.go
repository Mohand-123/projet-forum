package services

import (
	"errors"

	"playzone/repositories"
)

// BanUser bannit un compte
func BanUser(userID int) error {
	u, err := repositories.FindUserByID(userID)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}
	if u.Role == "admin" {
		return errors.New("impossible de bannir un admin")
	}
	return repositories.SetBanned(userID, true)
}

// UnbanUser debannit un compte
func UnbanUser(userID int) error {
	return repositories.SetBanned(userID, false)
}
