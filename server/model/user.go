package model

import (
	"context"
	"time"

	"circles/server/types"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/rs/xid"
)

// AllUsers retrieves all users and returns them as UserModels
func AllUsers(ctx context.Context, gs generalStore) ([]*UserModel, error) {
	users, err := gs.AllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userModels []*UserModel
	for _, user := range users {
		userModels = append(userModels, &UserModel{user, gs})
	}

	return userModels, nil
}

// ChatUsers retrieves all Chats for a User
func ChatUsers(ctx context.Context, gs generalStore, chid string) ([]*UserModel, error) {
	users, err := gs.AllUsersByChatMembership(ctx, chid)
	if err != nil {
		return nil, err
	}

	var userModels []*UserModel
	for _, u := range users {
		userModels = append(userModels, &UserModel{u, gs})
	}

	return userModels, nil
}

// CreateUser creates a new User with the given data and returns it as a UserModel
func CreateUser(ctx context.Context, gs generalStore, user *types.User) (*UserModel, error) {
	user.ID = xid.New().String()
	user.CreatedAt = time.Now().Format(time.RFC3339)

	if err := gs.CreateUser(ctx, user); err != nil {
		return nil, err

	} else if user, err = gs.FindUser(ctx, user.ID); err != nil {
		return nil, err
	}

	return &UserModel{user, gs}, nil
}

// DeleteUser deletes the user specified by ID and returns it as a UserModel
func DeleteUser(ctx context.Context, gs generalStore, id string) (*UserModel, error) {
	gs.BeginTransaction(ctx)
	defer gs.EndTransaction(ctx)

	user, err := gs.FindUser(ctx, id)
	if err != nil {
		return nil, err

	} else if err := gs.DeleteUser(ctx, id); err != nil {
		return nil, err

	} else if err := gs.DeleteMembershipsByUserID(ctx, id); err != nil {
		return nil, err
	}

	return &UserModel{user, gs}, nil
}

// FindUser retrieves the user specified by id and returns it as a UserModel
func FindUser(ctx context.Context, gs generalStore, id string) (*UserModel, error) {
	user, err := gs.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserModel{user, gs}, nil
}

// LoginUserInput ...
type LoginUserInput struct {
	PhoneNumber string
	DisplayName string
}

// LoginUser queries for a user using their PhoneNumber and DisplayName and returns that user as UserModel
func LoginUser(ctx context.Context, gs generalStore, input *LoginUserInput) (*UserModel, error) {
	user, err := gs.LoginUser(ctx, input.PhoneNumber, input.DisplayName)
	if err != nil {
		return nil, err
	}

	return &UserModel{user, gs}, nil
}

// UpdateUserInput ...
type UpdateUserInput struct {
	ID          string
	Name        *string
	PhoneNumber *string
	DisplayName *string
}

// UpdateUser updates the user specified by ID with the given data and returns it as a UserModel
func UpdateUser(ctx context.Context, gs generalStore, input *UpdateUserInput) (*UserModel, error) {
	user, err := gs.FindUser(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.PhoneNumber != nil {
		user.PhoneNumber = *input.PhoneNumber
	}

	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}

	if err = gs.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &UserModel{user, gs}, nil
}

// UserModel is the resolvable struct for the User struct
type UserModel struct {
	*types.User
	store generalStore
}

// ID field resolver
func (u *UserModel) ID() graphql.ID {
	return graphql.ID(u.User.ID)
}

// Name field resolver
func (u *UserModel) Name() string {
	return u.User.Name
}

// PhoneNumber field resolver
func (u *UserModel) PhoneNumber() string {
	return u.User.PhoneNumber
}

// DisplayName field resolver
func (u *UserModel) DisplayName() string {
	return u.User.DisplayName
}

// CreatedAt field resolver
func (u *UserModel) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, u.User.CreatedAt)
	return graphql.Time{Time: t}, err
}

// Chats field resolver
func (u *UserModel) Chats(ctx context.Context) ([]*ChatModel, error) {
	ctx = context.WithValue(ctx, types.Key("user_id"), u.User.ID)
	return UserChats(ctx, u.store, u.User.ID)
}
