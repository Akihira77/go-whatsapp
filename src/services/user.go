package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
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

func (us *UserService) Signup(ctx context.Context, data *types.Signup, image multipart.File) (*types.User, string, error) {
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		slog.Error("Hashing password",
			"error", err,
		)

		return nil, "", err
	}

	var fileByte []byte
	if image != nil {
		fileByte, err = io.ReadAll(image)
		if err != nil {
			slog.Error("Saving uploaded file",
				"error", err,
			)
			return nil, "", err
		}
	}

	user := types.User{
		ID:        ulid.Make().String(),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashedPassword,
		ImageUrl:  fileByte,
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

	go us.UpdateUserStatus(context.Background(), &user, types.ONLINE)
	user.Status = types.ONLINE

	return &user, token, err
}

func (us *UserService) Signin(ctx context.Context, data *types.Signin) (*types.User, string, error) {
	user, err := us.userRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		slog.Error("Finding user by email",
			"error", err,
		)

		return nil, "", fmt.Errorf("User not found")
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

		return nil, "", fmt.Errorf("Password invalid")
	}

	go us.UpdateUserStatus(context.Background(), user, types.ONLINE)
	user.Status = types.ONLINE

	return user, token, nil
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

func (us *UserService) FindUserByID(ctx context.Context, userId string) (*types.User, error) {
	user, err := us.userRepository.FindByID(ctx, userId)
	if err != nil {
		slog.Error("Finding user by id",
			"error", err,
		)
	}

	return user, err
}

func (us *UserService) FindGroupByID(ctx context.Context, groupId string) (*types.Group, error) {
	group, err := us.userRepository.FindGroupByID(ctx, groupId)
	if err != nil {
		slog.Error("Finding group by id",
			"error", err,
		)
	}

	return group, err
}

func (us *UserService) GetUserInfo(ctx context.Context, userId string) (*types.User, error) {
	user, err := us.userRepository.GetUserImage(ctx, userId)
	if err != nil {
		slog.Error("Get user image",
			"error", err,
		)
	}

	return user, err
}

func (us *UserService) GetMyContacts(ctx context.Context, userID, name string) ([]types.UserContact, error) {
	users, err := us.userRepository.FindMyContacts(ctx, userID, name)
	if err != nil {
		slog.Error("Finding my contacts",
			"error", err,
		)
	}

	return users, err
}

func (us *UserService) GetUsers(ctx context.Context, myUser *types.User, query *types.UserQuerySearch) ([]types.User, error) {
	users, err := us.userRepository.GetUsers(ctx, myUser, query)
	if err != nil {
		slog.Error("Get all users",
			"error", err,
		)
	}

	return users, err
}

func (us *UserService) UpdatePassword(ctx context.Context, user *types.User, data types.UpdatePassword) (*types.User, error) {
	isMatch := utils.CheckPasswordHash(data.OldPassword, user.Password)
	if !isMatch {
		slog.Error("Password did not match")
		return nil, fmt.Errorf("Password did not match")
	}

	hashedPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		slog.Error("Hashing password",
			"error", err,
		)
		return nil, err
	}

	user.Password = hashedPassword
	err = us.userRepository.Update(ctx, user)

	return user, err
}

func (us *UserService) UpdateUserProfile(ctx context.Context, myUser *types.User, data *types.UpdateUser, image multipart.File) (*types.User, error) {
	myUser.FirstName = data.FirstName
	myUser.LastName = data.LastName

	if image != nil {
		fileByte, err := io.ReadAll(image)
		if err != nil {
			slog.Error("Saving uploaded file",
				"error", err,
			)
			return nil, err
		}

		myUser.ImageUrl = fileByte
	}

	err := us.userRepository.Update(ctx, myUser)

	return myUser, err
}

func (us *UserService) AddContact(ctx context.Context, myUser *types.User, userID string) ([]types.UserContact, error) {
	contact := types.UserContact{
		UserOneID: myUser.ID,
		UserTwoID: userID,
		CreatedAt: time.Now(),
	}

	err := us.userRepository.AddContact(ctx, contact)
	if err != nil {
		slog.Error("Failed adding contact",
			"error", err,
		)
		return []types.UserContact{}, err
	}

	contacts, err := us.userRepository.FindMyContacts(ctx, myUser.ID, "")
	if err != nil {
		slog.Error("Retrieving user's contacts",
			"error", err,
		)
	}

	return contacts, err
}

func (us *UserService) RemoveContact(ctx context.Context, myUser *types.User, userID string) ([]types.UserContact, error) {
	contact := types.UserContact{
		UserOneID: myUser.ID,
		UserTwoID: userID,
		CreatedAt: time.Now(),
	}

	err := us.userRepository.RemoveContact(ctx, contact)
	if err != nil {
		slog.Error("Failed removing contact",
			"error", err,
		)
		return []types.UserContact{}, err
	}

	contacts, err := us.userRepository.FindMyContacts(ctx, myUser.ID, "")
	if err != nil {
		slog.Error("Retrieving user's contacts",
			"error", err,
		)
	}

	return contacts, err
}

func (us *UserService) FindGroups(ctx context.Context, userId string) ([]types.UserGroup, error) {
	groups, err := us.userRepository.FindGroups(ctx, userId)
	if err != nil {
		slog.Error("Retrieving your all groups",
			"error", err)
	}

	return groups, err
}

func (us *UserService) UpdateUserStatus(ctx context.Context, user *types.User, status types.UserStatus) (*types.User, error) {
	user.Status = status
	if status == types.OFFLINE {
		user.LastOnline = time.Now()
	}

	err := us.userRepository.Update(ctx, user)
	if err != nil {
		slog.Error("Updating your OFF/ON status",
			"error", err)
	}

	return user, err
}

func (us *UserService) CreateGroup(ctx context.Context, data types.CreateGroup, profile []byte, member []string) (*types.Group, error) {
	group, err := us.userRepository.CreateGroup(ctx, data, profile, member)
	if err != nil {
		slog.Error("Creating group",
			"error", err)
	}

	return group, err
}

func (us *UserService) EditGroup(ctx context.Context, group *types.Group, data types.EditGroup) (*types.Group, error) {
	if data.EditName {
		group.Name = data.Name
	}

	if data.EditDescription {
		group.Description = data.Description
	}

	group, err := us.userRepository.EditGroup(ctx, group)
	if err != nil {
		slog.Error("Editing group",
			"error", err)
	}

	return group, err
}
