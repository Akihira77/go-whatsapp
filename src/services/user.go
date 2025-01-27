package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Akihira77/go_whatsapp/src/repositories"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/Akihira77/go_whatsapp/src/utils"
	"github.com/oklog/ulid/v2"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Signup(ctx context.Context, data *types.Signup) (*types.User, string, error) {
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		slog.Error("Hashing password",
			"error", err,
		)

		return nil, "", err
	}

	user := types.User{
		ID:        ulid.Make().String(),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashedPassword,
		ImageUrl:  data.ImageUrl,
		CreatedAt: time.Now(),
	}

	err = us.userRepository.Create(ctx, &user)
	if err != nil {
		slog.Error("Creating user",
			"error", err,
		)
	}

	token, err := utils.GenerateJWT(&user)
	if err != nil {
		slog.Error("Generating JWT",
			"error", err,
		)

		return nil, "", err
	}

	return &user, token, err
}

func (us *UserService) Signin(ctx context.Context, data *types.Signin) (*types.User, string, error) {
	user, err := us.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		slog.Error("Finding user by email",
			"error", err,
		)

		return nil, "", err
	}

	isValid := utils.CheckPasswordHash(data.Password, user.Password)
	if !isValid {
		slog.Error("Password invalid")
		return nil, "", fmt.Errorf("Password invalid")
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		slog.Error("Generating JWT",
			"error", err,
		)

		return nil, "", err
	}

	return user, token, err
}

func (us *UserService) RefreshToken(ctx context.Context, tokenString string) (*types.User, string, error) {
	token, err := utils.VerifyingJWT(tokenString)
	if err != nil {
		slog.Error("Verifying jwt",
			"error", err,
		)

		return nil, "", err
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		slog.Error("Casting jwt claims",
			"claims", claims,
		)

		return nil, "", err
	}

	user, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		slog.Error("Finding user by email",
			"error", err,
		)

		return nil, "", err
	}

	newTokenString, err := utils.GenerateJWT(user)
	if err != nil {
		slog.Error("Generating JWT",
			"error", err,
		)

		return nil, "", err
	}

	return user, newTokenString, err
}

func (us *UserService) GetMyInfo(ctx context.Context, tokenString string) (*types.User, error) {
	token, err := utils.VerifyingJWT(tokenString)
	if err != nil {
		slog.Error("Verifying jwt",
			"error", err,
		)

		return nil, err
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		slog.Error("Casting jwt claims",
			"claims", claims,
		)

		return nil, err
	}

	user, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		slog.Error("Finding user by email",
			"error", err,
		)
	}

	return user, err
}

func (us *UserService) GetMyContacts(ctx context.Context, userID string) ([]types.UserContact, error) {
	users, err := us.userRepository.FindMyContacts(ctx, userID)
	if err != nil {
		slog.Error("Finding my contacts",
			"error", err,
		)
	}

	return users, err
}

func (us *UserService) GetUsers(ctx context.Context, myUser *types.User, query *types.UserQuerySearch) ([]types.UserContact, error) {
	users, err := us.userRepository.GetUsers(ctx, myUser, query)
	if err != nil {
		slog.Error("Get all users",
			"error", err,
		)
	}

	return users, err
}
