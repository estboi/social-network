package core

// Image Handler
func (c *Core) ImageUserProc(CurrentUserID int) ([]byte, error) {
	// Read image from db
	content, err := c.repo.ImageUserRead(CurrentUserID)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (c *Core) ImageUserCreateProc(CurrentUserID int, content []byte) error {
	// Record image to database
	if err := c.repo.CreateImage(CurrentUserID, content); err != nil {
		return err
	}
	return nil
}

