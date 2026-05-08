package user

import "golang.org/x/crypto/bcrypt"

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{
		Name:           input.Name,
		Email:          input.Email,
		Occupation:     input.Occupation,
		PasswordHash:   "",
		Role:           "user",
		AvatarFileName: "avatar.jpeg",
		Token:          "0",
	}
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(PasswordHash)
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}
