package core

import (
	"log"
	"social-network/core/entities"
)

// Users Handlers
func (c *Core) GetAllUsersProcessing(userId, From, To int) ([]entities.UsersVM, error) {
	users, err := c.repo.GetAllUsers(userId)
	if err != nil {
		log.Printf("Error getting users from database: %s\n", err)
		return []entities.UsersVM{}, err
	}
	return users, nil
}

func (c *Core) GetFollowedUsersProcessing(CurrentUserID int) ([]entities.UsersVM, error) {
	users, err := c.repo.GetFollowed(CurrentUserID)
	if err != nil {
		log.Println(err)
		return []entities.UsersVM{}, err
	}
	return users, nil
}

func (c *Core) GetFollowersProcessing(CurrentUserID, From, To int) ([]entities.UsersVM, error) {
	users, err := c.repo.GetFollowers(CurrentUserID)
	if err != nil {
		log.Println(err)
		return []entities.UsersVM{}, err
	}
	return users, nil
}

func (c *Core) GetUserProfileProcessing(CurrentUserID, UserID int) (entities.UserFullVM, error) {
	user, err := c.repo.GetProfile(CurrentUserID, UserID)
	if err != nil {
		log.Printf("Error getting user(%v) profile from db: %s", UserID, err)
		return entities.UserFullVM{}, err
	}
	return user, nil
}
