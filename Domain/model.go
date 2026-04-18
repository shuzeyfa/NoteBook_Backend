package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
	UserID    primitive.ObjectID `bson:"user_id" json:"-"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `json:"email" binding:"required,email"`
	Password string             `json:"password,omitempty" binding:"required,min=8"`
	Role     string             `json:"role"`
	GoogleID string             `bson:"google_id,omitempty" json:"-"` // new field for Google OAuth
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// Repository
type NoteRepository interface {
	GetAllNote(userID primitive.ObjectID) ([]Note, error)
	GetNoteByID(noteID, userID primitive.ObjectID) (Note, error)
	CreateNote(note Note, userID primitive.ObjectID) (Note, error)
	UpdateNote(note Note, userID primitive.ObjectID) (Note, error)
	DeleteNote(noteID, userID primitive.ObjectID) error
}

type UserRepository interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (User, error)
}

// AI related types
type GenerateAIRequest struct {
	Message string `json:"message" binding:"required"`
}

type GenerateAIResponse struct {
	Response string `json:"response"`
}
