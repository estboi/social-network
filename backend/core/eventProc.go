package core

import (
	"social-network/core/entities"
)

// GET
func (c *Core) EventGetAllProc(CurrentUserID int) ([]entities.EventVM, error) {
	events, err := c.repo.GetAllEvents(CurrentUserID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (c *Core) EventGetGroupProc(CurrentUserID, GroupID int) ([]entities.EventVM, error) {
	events, err := c.repo.GetGroupEvents(CurrentUserID, GroupID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// CREATE
func (c *Core) EventCreateProc(CurrentUserID int, EventDTO entities.EventDTO) ([]int, error) {
	users, err := c.repo.CreateEvent(CurrentUserID, EventDTO)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ACTIONS
func (c *Core) EventAcceptProc(CurrentUserID, EventID int) error {
	if err := c.repo.JoinEvent(CurrentUserID, EventID); err != nil {
		return err
	}
	return nil
}


func (c *Core) EventDenyProc(CurrentUserID, EventID int) error {
	if err := c.repo.DenyEvent(CurrentUserID, EventID); err != nil {
		return err
	}
	return nil
}