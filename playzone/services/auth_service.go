package services

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"

	"playzone/config"
	"playzone/models"
	"playzone/repositories"
)

// HashPassword genere un salt + hash SHA-512 du mdp
func HashPassword(password string) (string, string) {
	saltBytes := make([]byte, 16)
	rand.Read(saltBytes)
	salt := hex.EncodeToString(saltBytes)

	h := sha512.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil)), salt
}

// CheckPassword verifie qu'un mdp correspond au hash + salt stocke
func CheckPassword(password, hash, salt string) bool {
	h := sha512.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil)) == hash
}

// ValidatePassword applique les regles : 12 chars min, 1 maj, 1 special
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("le mot de passe doit faire au moins 12 caracteres")
	}
	hasUpper := false
	hasSpecial := false
	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			hasSpecial = true
		}
	}
	if !hasUpper {
		return errors.New("le mot de passe doit contenir au moins une majuscule")
	}
	if !hasSpecial {
		return errors.New("le mot de passe doit contenir au moins un caractere special")
	}
	return nil
}

// ValidateEmail check basique de format
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return re.MatchString(email)
}

// ValidateUsername check basique
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("le pseudo doit faire entre 3 et 20 caracteres")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !re.MatchString(username) {
		return errors.New("le pseudo ne doit contenir que des lettres, chiffres ou _")
	}
	return nil
}

// Register cree un nouveau compte apres validation
func Register(username, email, password string) (int64, error) {
	if err := ValidateUsername(username); err != nil {
		return 0, err
	}
	if !ValidateEmail(email) {
		return 0, errors.New("email invalide")
	}
	if err := ValidatePassword(password); err != nil {
		return 0, err
	}
	if repositories.ExistsByUsernameOrEmail(username, email) {
		return 0, errors.New("pseudo ou email deja utilise")
	}

	hash, salt := HashPassword(password)
	return repositories.CreateUser(username, email, hash, salt)
}

// Login verifie les credentials et retourne le user si OK
func Login(loginOrEmail, password string) (*models.User, error) {
	user, err := repositories.FindUserByLogin(loginOrEmail)
	if err != nil {
		return nil, errors.New("identifiants invalides")
	}
	if !CheckPassword(password, user.PasswordHash, user.Salt) {
		return nil, errors.New("identifiants invalides")
	}
	if user.IsBanned {
		return nil, errors.New("ce compte a ete banni")
	}
	return user, nil
}

// GenerateJWT cree un token pour un user
func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.JWTHours)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// ParseJWT valide un token et retourne les claims
func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("methode de signature invalide")
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide")
	}
	return claims, nil
}
