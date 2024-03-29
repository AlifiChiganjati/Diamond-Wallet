package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/RifaldyAldy/diamond-wallet/model"
	"github.com/RifaldyAldy/diamond-wallet/model/dto"
	"github.com/RifaldyAldy/diamond-wallet/repository"
	"github.com/RifaldyAldy/diamond-wallet/utils/common"
	encryption "github.com/RifaldyAldy/diamond-wallet/utils/encription"
)

type UserUseCase interface {
	CreateUser(payload dto.UserRequestDto) (model.User, error)
	LoginUser(in dto.LoginRequestDto) (dto.LoginResponseDto, error)
	FindById(id string) (model.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func (u *userUseCase) FindById(id string) (model.User, error) {
	user, err := u.repo.Get(id)
	if err != nil {
		return model.User{}, fmt.Errorf("user with ID %s not found", id)
	}

	return user, nil
}

func (u *userUseCase) CreateUser(payload dto.UserRequestDto) (model.User, error) {
	hashPassword, err := encryption.HashPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}
	newUser := model.User{
		Id:          payload.Id,
		Name:        payload.Name,
		Username:    payload.Username,
		Password:    hashPassword,
		Role:        payload.Role,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}
	user, err := u.repo.Create(newUser)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user: %v", err.Error())
	}
	return user, nil
}

func (u *userUseCase) LoginUser(in dto.LoginRequestDto) (dto.LoginResponseDto, error) {
	userData, err := u.repo.GetByUsername(in.Username)
	if err != nil {
		return dto.LoginResponseDto{}, err
	}
	isValid := encryption.CheckPasswordHash(in.Pass, userData.Password)
	if !isValid {
		return dto.LoginResponseDto{}, errors.New("1")
	}

	loginExpDuration := time.Duration(10) * time.Minute
	expiredAt := time.Now().Add(loginExpDuration).Unix()
	// TODO: tempel generate token jwt
	accessToken, err := common.GenerateTokenJwt(userData, expiredAt)
	if err != nil {
		return dto.LoginResponseDto{}, err
	}
	return dto.LoginResponseDto{
		AccessToken: accessToken,
		UserId:      userData.Id,
	}, nil
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}