package usecase

import (
	"errors"
	domain "taskmanagement/Domain"
	infrastructure "taskmanagement/Infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repo domain.UserRepository
}

func (u *UserUsecase) RegisterUser(req domain.RegisterRequest) (domain.User, error) {

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		12,
	)
	if err != nil {
		return domain.User{}, err
	}

	// check if the user already exist
	_, err = u.Repo.GetUserByEmail(req.Email)
	if err == nil {
		return domain.User{}, errors.New("user already exist")
	}

	user := domain.User{
		ID:       primitive.NewObjectID(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = u.Repo.CreateUser(user)

	return user, nil
}

func (u *UserUsecase) LoginUser(req domain.LoginRequest) (string, error) {

	// check if the user not registered
	user, err := u.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return "", errors.New("invalid credentials")
	}

	tokenString, err := infrastructure.GenerateJWT(user, req)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return tokenString, nil
}
