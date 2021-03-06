package store

import (
	"context"
	"fmt"

	"circles/server/types"
)

var (
	chatsTable = table("chats")

	allChatsByUserMembershipSQL = `
	SELECT * 
	FROM chats
	WHERE id IN (
		SELECT chat_id 
		FROM memberships
		WHERE user_id = $1
	)
	ORDER BY datetime(last_message_at) DESC`

	createChatSQL = `
	INSERT INTO chats (id, name, last_circle_name, last_message_at, created_at)
	VALUES (:id, :name, :last_circle_name, :last_message_at, :created_at)`

	updateChatSQL = `
	UPDATE chats 
	SET name=:name,
		last_circle_name=:last_circle_name,
		last_message_at=:last_message_at
	WHERE id=:id`
)

// AllChats finds all Chat entries in the db
func (gs *GeneralStore) AllChats(ctx context.Context) ([]*types.Chat, error) {
	var chats []*types.Chat
	if err := gs.findAllEntity(ctx, &chats, chatsTable); err != nil {
		return nil, fmt.Errorf("AllChats error: %v", err.Error())
	}

	return chats, nil
}

// AllChatsByUserMembership finds all Chat entries for a given User in the db
func (gs *GeneralStore) AllChatsByUserMembership(ctx context.Context, uid string) ([]*types.Chat, error) {
	var chats []*types.Chat
	if err := gs.getSet(ctx, &chats, allChatsByUserMembershipSQL, uid); err != nil {
		return nil, fmt.Errorf("AllChatsByUserMembership error: %v", err.Error())
	}

	return chats, nil
}

// CreateChat creates a Chat entry in the db
func (gs *GeneralStore) CreateChat(ctx context.Context, chat *types.Chat, userIDs []string) error {
	if err := gs.createEntity(ctx, chat, createChatSQL); err != nil {
		return fmt.Errorf("CreateChat error: %v", err.Error())
	}

	return nil
}

// DeleteChat deletes a Chat entry in the db, and any foreign related entries
func (gs *GeneralStore) DeleteChat(ctx context.Context, id string) error {
	if err := gs.deleteEntity(ctx, chatsTable, id); err != nil {
		return fmt.Errorf("DeleteChat error: %v", err.Error())
	}

	return nil
}

// FindChat finds a Chat entry in the db
func (gs *GeneralStore) FindChat(ctx context.Context, id string) (*types.Chat, error) {
	chat := &types.Chat{}
	if err := gs.findEntity(ctx, chat, chatsTable, id); err != nil {
		return nil, fmt.Errorf("FindChat error: %v", err.Error())
	}

	return chat, nil
}

// UpdateChat updates a Chat entry in the db
func (gs *GeneralStore) UpdateChat(ctx context.Context, chat *types.Chat) error {
	if err := gs.updateEntity(ctx, chat, updateChatSQL); err != nil {
		return fmt.Errorf("UpdateChat error: %v", err.Error())
	}

	return nil
}
