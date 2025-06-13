package core

import (
	"fmt"
	"social-network/core/entities"
)

// Chats Handlers

func (c *Core) SendMessage(msg entities.Message) error {
	fmt.Printf("SendMessage method called with Id: %d\n", msg.SenderId)
	err := c.repo.SaveMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) GetChatMessages(CurrentUserID, chatID int) ([]entities.Message, error) {
	fmt.Printf("ChatsHandlerGetGroup method called with Id: %d\n", chatID)
	history := c.repo.GetOneUserMessages(CurrentUserID, chatID)
	return history, nil
}

func (c *Core) ChatGetGroupMessagesProc(CurrentUserID, groupID int) ([]entities.Message, error) {
	fmt.Printf("ChatsHandlerGetGroup method called with Id: %d\n", CurrentUserID)
	history := c.repo.GetOneGroupMessages(CurrentUserID, groupID)
	return history, nil
}

func (c *Core) ChatGetAllMessagesProc(CurrentUserID int) ([]entities.Message, error) {
	fmt.Printf("ChatsHandlerGetAll method called with Id: %d\n", CurrentUserID)
	history := c.repo.GetUserChats(CurrentUserID)
	groupHistory := c.repo.GetGroupChats(CurrentUserID)
	history = append(history, groupHistory...)
	return history, nil
}
