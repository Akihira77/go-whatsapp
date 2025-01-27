package repositories

import (
	"context"

	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/Akihira77/go_whatsapp/src/types"
)

type UserRepository struct {
	store *store.Store
}

func NewUserRepository(store *store.Store) *UserRepository {
	return &UserRepository{
		store: store,
	}
}

func (ur *UserRepository) FindByID(ctx context.Context, id string) (*types.User, error) {
	var u types.User

	res := ur.
		store.
		DB.
		Model(&types.User{}).
		WithContext(ctx).
		Where("id = ?", id).
		First(&u)

	return &u, res.Error
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	var u types.User

	res := ur.
		store.
		DB.
		Model(&types.User{}).
		WithContext(ctx).
		Where("email = ?", email).
		First(&u)

	return &u, res.Error
}

func (ur *UserRepository) SearchByName(ctx context.Context, name string) ([]types.User, error) {
	var users []types.User

	res := ur.
		store.
		DB.
		Model(&types.User{}).
		WithContext(ctx).
		Where("first_name || ' ' || last_name LIKE ?", "%"+name+"%").
		Find(&users)

	return users, res.Error
}

func (ur *UserRepository) Create(ctx context.Context, data *types.User) error {
	res := ur.
		store.
		DB.
		Model(&types.User{}).
		WithContext(ctx).
		Create(&data)

	return res.Error
}

func (ur *UserRepository) Update(ctx context.Context, data *types.User) error {
	res := ur.
		store.
		DB.
		WithContext(ctx).
		Save(&data)

	return res.Error
}

func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	res := ur.
		store.
		DB.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&types.User{})

	return res.Error
}

func (ur *UserRepository) FindMyContacts(ctx context.Context, userID string) ([]types.UserContact, error) {
	var users []types.UserContact

	res := ur.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Model(&types.UserContact{}).
		Preload("UserOne", "id <> ?", userID).
		Preload("UserTwo", "id <> ?", userID).
		Where("user_one_id = ? OR user_two_id = ?", userID, userID).
		Find(&users)

	return users, res.Error
}
