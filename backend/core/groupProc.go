package core

import (
	"errors"
	"social-network/core/entities"
)

// Get
func (c *Core) GroupsAllProc(userID int) ([]entities.GroupVM, error) {
	// Call the DB method
	groups, err := c.repo.GetAllGroups(userID)
	if err != nil {
		return []entities.GroupVM{}, err
	}

	return groups, nil
}

func (c *Core) GroupsHandlerGetConnected(CurrentUserID int) ([]entities.GroupVM, error) {
	groups, err := c.repo.GetConnectedGroups(CurrentUserID)
	if err != nil {
		return []entities.GroupVM{}, err
	}

	return groups, nil
}

func (c *Core) GroupsHandlerGetCreated(CurrentUserID int) ([]entities.GroupVM, error) {
	groups, err := c.repo.GetCreatedGroups(CurrentUserID)
	if err != nil {
		return []entities.GroupVM{}, err
	}

	return groups, nil
}

func (c *Core) GroupsProfileProc(UserId, GroupID int) (entities.GroupVM, error) {
	groups, err := c.repo.GroupsProfileRead(UserId, GroupID)
	if err != nil {
		return entities.GroupVM{}, err
	}
	return groups, nil
}

func (c *Core) GroupsGetNotMembersProc(groupId int) ([]entities.UsersVM, error) {
	users, err := c.repo.GroupsGetNotMembers(groupId)
	if err != nil {
		return []entities.UsersVM{}, err
	}

	return users, nil
}

func (c *Core) GroupsGetRequestedProc(userId, groupId int) ([]entities.UsersVM, error) {

	users, err := c.repo.GroupsGetRequested(userId, groupId)
	if err != nil {
		return []entities.UsersVM{}, err
	}

	return users, nil
}

// Create
func (c *Core) GroupsCreateProc(CurrentUserID int, dto entities.GroupDTO) (int, error) {
	// Step 1: Check input validity
	switch {
	case len(dto.GroupName) == 0:
		return -1, errors.New("group Name should be present")
	case len(dto.GroupName) > 25:
		return -1, errors.New("group Name max length 25 chars")
	case len(dto.GroupAbout) == 0:
		return -1, errors.New("group Description should be present")
	case len(dto.GroupAbout) > 250:
		return -1, errors.New("group Description max length 25 chars")
	}

	// Step 2: call repo method
	groupId, err := c.repo.CreateGroup(CurrentUserID, dto)
	if err != nil {
		return -1, err
	}
	return groupId, nil
}

// IMAGES
func (c *Core) GroupsImageCreateProc(groupId int, content []byte) error {
	// Step 1: Call repo method
	if err := c.repo.GroupsRecordImage(groupId, content); err != nil {
		return err
	}
	return nil
}

func (c *Core) GroupsImageProc(groupId int) ([]byte, error) {
	content, err := c.repo.GroupsImageRead(groupId)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// ACTIONS
func (c *Core) GroupsRequestProc(CurrentUserID, GroupID int) (int, error) {
	creatorID, err := c.repo.GroupRequest(CurrentUserID, GroupID)
	if err != nil {
		return -1, err
	}
	return creatorID, nil
}

func (c *Core) GroupsInviteProc(InvitedUserId, GroupID int) error {
	if err := c.repo.GroupInvite(InvitedUserId, GroupID); err != nil {
		return err
	}
	return nil
}

func (c *Core) GroupsAcceptProc(RequestingUserId, GroupID int) error {
	if err := c.repo.GroupAccept(RequestingUserId, GroupID); err != nil {
		return err
	}
	return nil
}

func (c *Core) GroupsDenyProc(RequestingUserId, GroupID int) error {
	if err := c.repo.GroupDeny(RequestingUserId, GroupID); err != nil {
		return err
	}
	return nil
}
