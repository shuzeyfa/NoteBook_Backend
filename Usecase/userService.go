package usecase

import (
	"errors"
	domain "taskmanagement/Domain"
	infrastructure "taskmanagement/Infrastructure"
)

type UserUsecase struct {
	Repo domain.UserRepository
}

func (u *UserUsecase) RegisterUser(req domain.RegisterRequest) (domain.User, error) {
	user := domain.User{
		Email: req.Email,
		Role:  "user",
	}

	hashed, _ := infrastructure.HashPassword(req.Password)
	user.Password = hashed

	err := u.Repo.CreateUser(user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserUsecase) LoginUser(req domain.LoginRequest) (string, error) {
	user, err := u.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return infrastructure.GenerateJWT(user, req)
}
