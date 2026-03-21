package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	DueDate     string             `bson:"due_date" json:"due_date"`
	Status      string             `bson:"status" json:"status"`
	UserId      primitive.ObjectID `bson:"user_id,omitempty"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email" binding:"required,email"`
	Password string             `json:"password,omitempty" binding:"required,min=8"`
	Role     string             `json:"role"`
}

type Claims struct {
	UserID primitive.ObjectID `bson:"_id,omitempty"`
	Email  string             `json:"email"`
	Role   string             `json:"role"`
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

type UserRepository interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (User, error)
}
