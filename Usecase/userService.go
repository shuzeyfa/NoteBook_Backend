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

// Google signin method
func (u *UserUsecase) GoogleLogin(req domain.GoogleLoginRequest) (string, error) {
	payload, err := infrastructure.VerifyGoogleIDToken(req.IDToken)
	if err != nil {
		return "", errors.New("invalid google token")
	}

	email := payload.Claims["email"].(string)
	googleID := payload.Claims["sub"].(string) // unique Google user ID

	// Try to find existing user by email (most common case)
	user, err := u.Repo.GetUserByEmail(email)
	if err != nil {
		// User doesn't exist → create new Google user
		newUser := domain.User{
			Email:    email,
			Role:     "user",
			GoogleID: googleID,
			// Password left empty - they signed up via Google
		}
		if err := u.Repo.CreateUser(newUser); err != nil {
			return "", err
		}
		user = newUser // now we have the user
	} else {
		// User exists → optional: update GoogleID if it was missing
		if user.GoogleID == "" {
			user.GoogleID = googleID
			// You could call an UpdateUser method here if you have one
		}
	}

	// Generate your JWT exactly like normal login
	return infrastructure.GenerateJWTForAuthenticatedUser(user)
}
