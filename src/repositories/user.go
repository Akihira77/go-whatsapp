package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
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

func (ur *UserRepository) GetUserImage(ctx context.Context, id string) (*types.User, error) {
	var u types.User

	res := ur.
		store.
		DB.
		Model(&types.User{}).
		WithContext(ctx).
		Select("image_url").
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

func (ur *UserRepository) FindMyContacts(ctx context.Context, userID, name string) ([]types.UserContact, error) {
	var users []types.UserContact

	res := ur.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Model(&types.UserContact{}).
		Where("user_one_id = ? OR user_two_id = ?", userID, userID).
		Preload("UserOne", "id <> ? AND (first_name || ' ' || last_name) LIKE ?", userID, "%"+name+"%").
		Preload("UserTwo", "id <> ? AND (first_name || ' ' || last_name) LIKE ?", userID, "%"+name+"%").
		// Preload("UserOne").
		// Preload("UserTwo").
		Find(&users)

	return users, res.Error
}

func (ur *UserRepository) GetUsers(ctx context.Context, myUser *types.User, query *types.UserQuerySearch) ([]types.User, error) {
	var users []types.User

	res := ur.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Model(&types.User{}).
		Where("id <> ? AND (first_name || ' ' || last_name) LIKE ?", myUser.ID, "%"+query.Search+"%").
		Offset((query.Page - 1) * query.Size).
		Limit(query.Size).
		Order("(first_name || ' ' || last_name) ASC").
		Find(&users)

	return users, res.Error
}

func (ur *UserRepository) AddContact(ctx context.Context, data types.UserContact) error {
	res := ur.
		store.
		DB.
		Debug().
		Model(&types.UserContact{}).
		Create(&data)

	return res.Error
}

func (ur *UserRepository) RemoveContact(ctx context.Context, data types.UserContact) error {
	res := ur.
		store.
		DB.
		Debug().
		Model(&types.UserContact{}).
		Where("user_one_id = ? AND user_two_id = ?", data.UserOneID, data.UserTwoID).
		Delete(&types.UserContact{})

	return res.Error
}

func (ur *UserRepository) FindGroups(ctx context.Context, userId string) ([]types.UserGroup, error) {
	var groups []types.UserGroup

	res := ur.
		store.
		DB.
		Debug().
		Model(&types.UserGroup{}).
		Preload("Group").
		Where("user_id = ?", userId).
		Find(&groups)

	return groups, res.Error
}

func (ur *UserRepository) CreateGroup(ctx context.Context, data types.CreateGroup, groupProfile []byte, member []string) (*types.Group, error) {
	tx := ur.store.DB.Begin(&sql.TxOptions{}).Debug()

	group := types.Group{
		ID:           ulid.Make().String(),
		Name:         data.Name,
		UserCount:    len(member),
		CreatorID:    data.Creator.ID,
		GroupProfile: groupProfile,
		CreatedAt:    time.Now(),
	}

	res := tx.
		Model(&types.Group{}).
		Create(&group)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	userGroups := make([]types.UserGroup, 0)
	for _, userId := range member {
		userGroups = append(userGroups, types.UserGroup{
			GroupID: group.ID,
			UserID:  userId,
		})
	}

	res = tx.
		Model(&types.UserGroup{}).
		Create(&userGroups)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	res = tx.Commit()
	group.Member = userGroups

	return &group, res.Error
}

func (ur *UserRepository) EditGroup(ctx context.Context, group *types.Group) (*types.Group, error) {
	res := ur.
		store.
		DB.
		Debug().
		WithContext(ctx).
		Save(group)

	return group, res.Error
}

func (ur *UserRepository) FindGroupByID(ctx context.Context, id string) (*types.Group, error) {
	var u types.Group

	res := ur.
		store.
		DB.
		Debug().
		Model(&types.Group{}).
		WithContext(ctx).
		Preload("Creator").
		Preload("Member", func(tx *gorm.DB) *gorm.DB {
			return tx.Limit(10)
		}).
		Preload("Member.User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "first_name", "last_name", "email")
		}).
		Where("id = ?", id).
		First(&u)

	return &u, res.Error
}

func (ur *UserRepository) GetGroupMembers(ctx context.Context, groupId string) ([]types.UserGroup, error) {
	var members []types.UserGroup

	res := ur.
		store.
		DB.
		Debug().
		Model(&types.UserGroup{}).
		WithContext(ctx).
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "first_name", "last_name", "email")
		}).
		Where("group_id = ?", groupId).
		Find(&members)

	return members, res.Error
}
