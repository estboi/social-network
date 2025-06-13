package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/core/entities"
)

// BASE USERS INFO
func (d Database) CreateUser(user entities.RegisterDTO) (int, error) {
	query := `INSERT INTO Users(FirstName, LastName, NickName, Email, PasswordHash, DateOfBirth, About) VALUES(?, ?, ?, ?, ?, ?, ?)`
	statement, err := d.db.Prepare(query)
	if err != nil {
		return -1, err
	}
	result, err := statement.Exec(user.FirstName, user.LastName, user.NickName, user.Email, user.Pass, user.Date, user.About)
	if err != nil {
		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(userID), nil
}

func (d Database) GetAllUsers(CurrentUserID int) ([]entities.UsersVM, error) {
	var users []entities.UsersVM

	rows, err := d.db.Query(`
	SELECT 
	ID, 
	COALESCE(Users.NickName, Users.FirstName || ' ' || Users.LastName) AS UserName 
	FROM Users
	WHERE ID != ?`, CurrentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u entities.UsersVM
		err = rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return nil, err
		}

		// Set IsFollower and IsFollow based on relationships
		u.IsFollower = d.IsUserFollowing(CurrentUserID, u.ID)
		u.IsFollow = d.IsUserFollowed(CurrentUserID, u.ID)
		if u.IsFollow {
			u.IsPending = d.IsUserPending(u.ID, CurrentUserID)
		}

		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (d Database) GetUserName(CurrentUserID int) (entities.NavbarVM, error) {

	var vm entities.NavbarVM

	// prepare query
	query := `SELECT CASE
    WHEN NickName IS NOT NULL THEN NickName
    ELSE FirstName || ' ' || LastName
    END AS DisplayedName
	FROM Users
	WHERE ID = ?;
	`

	// execute the query
	row := d.db.QueryRow(query, CurrentUserID)

	// store the result in entities.NavbarVM struct
	err := row.Scan(&vm.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return vm, fmt.Errorf("no such user exists with the ID: %v", CurrentUserID)
		}
		return vm, err
	}

	vm.UserID = CurrentUserID

	return vm, nil
}

func (d Database) GetProfile(CurrentUserID, UserID int) (entities.UserFullVM, error) {
	var user entities.UserFullVM
	query := `
		SELECT 
		ID, 
		FirstName || ' ' || LastName as Name, 
		COALESCE(NickName, '') as NickName, 
		Email, 
		DateOfBirth, 
		COALESCE(About, '') as About, 
		Privacy 
		FROM Users 
		WHERE ID = ?`

	row := d.db.QueryRow(query, UserID)
	err := row.Scan(&user.Id, &user.Name, &user.Nickname, &user.Email, &user.Date, &user.About, &user.Privacy)
	if err != nil {
		return entities.UserFullVM{}, err
	}

	if UserID != CurrentUserID && user.Privacy == "Private" && (!d.IsUserFollowing(CurrentUserID, UserID) || d.IsUserPending(CurrentUserID, UserID)) {
		return entities.UserFullVM{}, errors.New("no access")
	}
	return user, nil
}

func (d Database) GetPasswordHash(userId int) (string, error) {
	query := `SELECT PasswordHash FROM Users WHERE ID=?`

	var passwordHash string
	err := d.db.QueryRow(query, userId).Scan(&passwordHash)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %v", err)
	}

	return passwordHash, nil
}

func (d Database) LoginUser(User entities.LoginDTO) (entities.LoginValues, error) {
	var loginValues entities.LoginValues

	rows, err := d.db.Query(`SELECT ID FROM Users WHERE Email=? OR NickName=?`, User.Login, User.Login)
	if err != nil {
		return loginValues, err
	}
	defer rows.Close()

	// Read the results
	for rows.Next() {
		err := rows.Scan(&loginValues.UserId)
		if err != nil {
			return loginValues, err
		}
	}
	// Get PasswordHash
	if loginValues.UserId != 0 {
		loginValues.HashPassword, err = d.GetPasswordHash(loginValues.UserId)
		if err != nil {
			return loginValues, err
		}
	}
	return loginValues, nil
}

func (d Database) ModifyUser(UserID int) (entities.Status, error) {
	var status entities.Status

	// Check the current privacy setting
	currentPrivacy, err := d.GetUserPrivacy(UserID)
	if err != nil {
		return status, err
	}

	// Toggle privacy
	newPrivacy := togglePrivacy(currentPrivacy)

	// Update the user's privacy in the database
	updateQuery := "UPDATE Users SET Privacy = ? WHERE ID = ?"
	_, err = d.db.Exec(updateQuery, newPrivacy, UserID)
	if err != nil {
		return status, err
	}

	status.Message = "Privacy updated successfully"
	return status, nil
}

// Function to toggle privacy between "Public" and "Private"
func togglePrivacy(currentPrivacy string) string {
	if currentPrivacy == "Public" {
		return "Private"
	} else if currentPrivacy == "Private" {
		return "Public"
	}
	// Handle other cases if needed
	return currentPrivacy
}

// Function to get the current privacy setting for a user
func (d Database) GetUserPrivacy(UserID int) (string, error) {
	var privacy string
	err := d.db.QueryRow("SELECT Privacy FROM Users WHERE ID = ?", UserID).Scan(&privacy)
	if err != nil {
		return "", err
	}
	return privacy, nil
}
