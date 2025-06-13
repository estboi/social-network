package helpers

import (
	"errors"
	"regexp"
	"social-network/core/entities"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidRegistrationInput(registerDTO entities.RegisterDTO) error {
	switch {
	case len(registerDTO.FirstName) == 0:
		return errors.New("firstname field should be filled")
	case len(registerDTO.FirstName) > 50:
		return errors.New("firstname max length is 50 chars")
	case len(registerDTO.LastName) == 0:
		return errors.New("lastname field should be filled")
	case len(registerDTO.LastName) > 50:
		return errors.New("lastname max length is 50 chars")
	case len(registerDTO.Email) == 0:
		return errors.New("email field should be filled")
	case !IsValidEmail(registerDTO.Email):
		return errors.New("email input is not correct")
	case len(registerDTO.Pass) == 0:
		return errors.New("password field should be filled")
	case len(registerDTO.Pass) < 6:
		return errors.New("password should contain more than 6 chars")
	case len(registerDTO.Date) == 0:
		return errors.New("pass field should be filled")
	case registerDTO.NickName != nil && len(*registerDTO.NickName) > 25:
		return errors.New("nickame max length is 25 chars")
	case len(registerDTO.About) > 500:
		return errors.New("about field length is 500 chars")
	}
	return nil
}
