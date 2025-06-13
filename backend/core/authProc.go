package core

import (
	"errors"
	helpers "social-network/core/common"
	"social-network/core/entities"
	crypt "social-network/crypting"
	"strings"
)

// Authentication Handlers
func (c *Core) LoginProcess(loginDTO entities.LoginDTO) (int, error) {
	// Step 1: validate input
	switch {
	case len(loginDTO.Login) == 0:
		return -1, errors.New("login field should be filled")
	case strings.Contains(loginDTO.Login, "@"):
		if !helpers.IsValidEmail(loginDTO.Login) {
			return -1, errors.New("email input is not correct")
		}
	case len(loginDTO.Pass) == 0:
		return -1, errors.New("pass field should be filled")
	}

	// Step 2: Check if this user exists
	loginValues, err := c.repo.LoginUser(loginDTO)
	if err != nil {
		return -1, errors.New("this user doesn't exists")
	}
	if err := crypt.CheckHash(loginDTO.Pass, loginValues.HashPassword); err != nil {
		return -1, errors.New("incorrect password")
	}

	return loginValues.UserId, nil
}

func (c *Core) RegisterProc(registerDTO entities.RegisterDTO) (int, error) {
	var err error
	// Step 1: Check for correct input
	if err = helpers.IsValidRegistrationInput(registerDTO); err != nil {
		return -1, err
	}

	// Step 2: Encrpyt the password
	registerDTO.Pass, err = crypt.HashPassword(registerDTO.Pass)
	if err != nil {
		return -1, errors.New("incorrect password")
	}

	// Step 3: record new user to DB
	userId, err := c.repo.CreateUser(registerDTO)
	if err != nil {
		return -1, err
	}
	return userId, nil
}
