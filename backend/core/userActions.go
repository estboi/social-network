package core

import "social-network/core/entities"

func (c *Core) SubscribeProcessing(CurrentUserID, UserID int) (entities.Status, error) {
	c.repo.SubscribeOnUser(CurrentUserID, UserID)
	return entities.Status{Message: "X action completed successfully )"}, nil
}

func (c *Core) UnsubscribeProcessing(CurrentUserID, UserID int) (entities.Status, error) {
	c.repo.UnsubscribeOnUser(CurrentUserID, UserID)
	return entities.Status{Message: "X action completed successfully )"}, nil
}

func (c *Core) ModifyUserProcessing(CurrentUserID int) (entities.Status, error) {
	c.repo.ModifyUser(CurrentUserID)
	return entities.Status{Message: "X action completed successfully )"}, nil
}
