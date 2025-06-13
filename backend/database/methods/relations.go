package sqlite

import (
	"database/sql"
	"log"
	"social-network/core/entities"
)

// USERS RELATIONS

func (d Database) SubscribeOnUser(CurrentUserID, UserID int) (entities.Status, error) {
	status := entities.Status{}

	currentPrivacy, err := d.GetUserPrivacy(UserID)
	if err != nil {
		log.Println(err)
	}

	if d.IsUserPending(UserID, CurrentUserID) {
		query := "UPDATE UserSubscriptionRelations SET CurrentStatus = ? WHERE UserID = ? AND FollowerID = ?"
		d.db.Exec(query, "Subscribed", CurrentUserID, UserID)
	}

	subscriptionStatus := "Subscribed"
	if currentPrivacy == "Private" {
		subscriptionStatus = "Pending"
	}

	statement, err := d.db.Prepare("INSERT INTO UserSubscriptionRelations (UserID, FollowerID, CurrentStatus) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(UserID, CurrentUserID, subscriptionStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return status, nil
		}
		return status, err
	}
	status.Message = "Successfully subscribed"
	return status, nil
}

func (d Database) UnsubscribeOnUser(CurrentUserID, UserID int) (entities.Status, error) {
	// Check if a subscription exists before attempting to remove it.
	subscriptionQuery := "SELECT ID FROM UserSubscriptionRelations WHERE UserID = ? AND FollowerID = ?;"
	var subscriptionID int
	err := d.db.QueryRow(subscriptionQuery, UserID, CurrentUserID).Scan(&subscriptionID)
	if err != nil {
		err = d.db.QueryRow(subscriptionQuery, CurrentUserID, UserID).Scan(&subscriptionID)
		if err == sql.ErrNoRows {
			// The subscription does not exist. Return an appropriate status.
			return entities.Status{Message: "Subscription does not exist."}, nil
		}
	}

	// If the subscription exists, delete it.
	deleteQuery := "DELETE FROM UserSubscriptionRelations WHERE ID = ?;"
	_, err = d.db.Exec(deleteQuery, subscriptionID)
	if err != nil {
		return entities.Status{}, err
	}

	return entities.Status{Message: "Unsubscribed successfully."}, nil
}

// is Other user following current user
func (d Database) IsUserFollowing(CurrentUserID, UserID int) bool {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM UserSubscriptionRelations WHERE UserID = ? AND FollowerID = ?", UserID, CurrentUserID).Scan(&count)
	if err != nil {
		// Handle the error, e.g., log it or return false
		return false
	}
	return count > 0
}

// is other user following current user with private profile
func (d Database) IsUserPending(CurrentUserID, UserID int) bool {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM UserSubscriptionRelations WHERE FollowerID = ? AND UserID = ? AND CurrentStatus = 'Pending'", CurrentUserID, UserID).Scan(&count)
	if err != nil {
		// Handle the error, e.g., log it or return false
		return false
	}
	return count > 0
}

// is current user following other user
func (d Database) IsUserFollowed(CurrentUserID, UserID int) bool {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM UserSubscriptionRelations WHERE UserID = ? AND FollowerID = ?", CurrentUserID, UserID).Scan(&count)
	if err != nil {
		// Handle the error, e.g., log it or return false
		return false
	}
	return count > 0
}

func (d Database) GetFollowed(CurrentUserID int) ([]entities.UsersVM, error) {
	var users []entities.UsersVM
	query := `SELECT 
	ID,
	COALESCE(U.NickName, U.FirstName || ' ' || U.LastName) AS UserName 
	FROM Users AS U
	`

	rows, err := d.db.Query(query)
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
		if d.IsUserFollowing(CurrentUserID, u.ID) {
			u.IsFollower = true
			users = append(users, u)
		}

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (d Database) GetFollowers(CurrentUserID int) ([]entities.UsersVM, error) {
	var users []entities.UsersVM

	rows, err := d.db.Query("SELECT ID, COALESCE(Users.NickName, Users.FirstName || ' ' || Users.LastName) AS UserName FROM Users")
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
		if d.IsUserFollowed(CurrentUserID, u.ID) {
			if d.IsUserFollowing(CurrentUserID, u.ID) {
				u.IsFollower = true
			}
			u.IsFollow = true
			users = append(users, u)
		}

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
