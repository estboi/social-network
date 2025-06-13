package core

import (
	"social-network/core/entities"
	"social-network/core/interfaces"
)

type Core struct {
	repo interfaces.Repository
}

func NewCore(repo interfaces.Repository) *Core {
	return &Core{
		repo: repo,
	}
}

// Navbar Handler
func (c *Core) NavbarProc(UserID int) (entities.NavbarVM, error) {
	// Step 1: get the values from DB
	navbarValues, err := c.repo.GetUserName(UserID)
	if err != nil {
		return navbarValues, err
	}

	return navbarValues, nil
}

// Notifications Handler
func (c *Core) NotificationsHandlerGet(CurrentUserID int) ([]entities.NotificationVM, error) {
	notifications, err := c.repo.GetNotifications(CurrentUserID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (c *Core) NotificationRecord(notif entities.NotificationVM) error {
	if err := c.repo.NotificationRecord(notif); err != nil {
		return err
	}
	return nil
}
