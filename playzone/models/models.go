package models

// User represente un utilisateur de la plateforme
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	Salt         string
	Role         string // "user" ou "admin"
	IsBanned     bool
	CreatedAt    string
}

// Category represente une categorie de jeux
type Category struct {
	ID   int
	Name string
	Slug string
}

// Tag represente un tag attache a un fil
type Tag struct {
	ID    int
	Name  string
	Color string
}

// Thread represente un fil de discussion
type Thread struct {
	ID         int
	Title      string
	Content    string
	AuthorID   int
	CategoryID int
	State      string // "ouvert", "ferme", "archive"
	CreatedAt  string
	UpdatedAt  string

	// Champs calcules pour les vues
	AuthorName   string
	CategoryName string
	CategorySlug string
	Tags         []Tag
	MessageCount int
}

// Message represente une reponse dans un fil
type Message struct {
	ID        int
	ThreadID  int
	AuthorID  int
	Content   string
	CreatedAt string
	UpdatedAt string

	// Champs calcules
	AuthorName string
	Score      int    // likes - dislikes
	UserReact  string // reaction de l'user courant : "like", "dislike", ou ""
}

// Reaction represente un like/dislike sur un message
type Reaction struct {
	ID        int
	UserID    int
	MessageID int
	Type      string // "like" ou "dislike"
}
