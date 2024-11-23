package auth

import (
	"errors"
	"order-api/internal/user"
	"order-api/pkg/session"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(phone, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	sessionId, err := session.GenerateSessionID()

	if err != nil {
		return "", err
	}

	user := &user.User{
		Phone:     phone,
		Name:      name,
		SessionId: sessionId,
	}

	_, err = service.UserRepository.Create(user)

	if err != nil {
		return "", err
	}

	return user.SessionId, nil

}

func (service *AuthService) Login(phone string) (string, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)
	if existedUser == nil {
		return "", errors.New(WrongCredentials)
	}

	sessionId, err := session.GenerateSessionID()
	if err != nil {
		return "", err
	}

	user := &user.User{
		Phone:     existedUser.Phone,
		Name:      existedUser.Name,
		SessionId: sessionId,
	}

	updatedUser, err := service.UserRepository.Update(user)

	if err != nil {
		return "", err
	}

	return updatedUser.SessionId, nil
}
