package sqlite

import (
	"database/sql"
	"fmt"
)

// IMAGES
func (d Database) IsAccess2Image(CurrentUserID int, Src string, ID int) (bool, error) {
	// Query if a post with the given source exists from the current user
	stmt, err := d.db.Prepare("SELECT Img FROM Posts WHERE CreatorID = ? AND ID = ?")
	if err != nil {
		return false, fmt.Errorf("prepare statement error: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(CurrentUserID, ID)

	var img string
	err = row.Scan(&img)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("queryRow error: %v", err)
	default:
		// If post exists, compare the source and return true if it matches, false otherwise
		return img == Src, nil
	}
}
func (d Database) ImageUserRead(UserId int) ([]byte, error) {
	var Img []byte
	// Prepare the sql statement
	stmt, err := d.db.Prepare("SELECT Img FROM Users WHERE ID=?")
	if err != nil {
		return nil, fmt.Errorf("cannot prepare the statement: %v", err)
	}
	defer stmt.Close()

	// Execute the sql statement
	err = stmt.QueryRow(UserId).Scan(&Img)
	if err != nil {
		return nil, fmt.Errorf("cannot execute the sql statement: %v", err)
	}

	return Img, nil
}
func (d Database) CreateImage(CurrentUserID int, Content []byte) error {
	stmt, err := d.db.Prepare("UPDATE Users SET Img = ? WHERE ID = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(Content, CurrentUserID)
	if err != nil {
		return err
	}

	return nil
}

func (d Database) RecordPostImage(PostId int, Content []byte) error {
	stmt, err := d.db.Prepare("UPDATE Posts SET Img = ? WHERE ID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Content, PostId)

	if err != nil {
		return nil
	}

	return nil
}

func (d Database) ReadPostImage(PostId int) ([]byte, error) {
	var Img []byte
	// Prepare the sql statement
	stmt, err := d.db.Prepare("SELECT Img FROM Posts WHERE ID=?")
	if err != nil {
		return nil, fmt.Errorf("cannot prepare the statement: %v", err)
	}
	defer stmt.Close()

	// Execute the sql statement
	err = stmt.QueryRow(PostId).Scan(&Img)
	if err != nil {
		return nil, fmt.Errorf("cannot execute the sql statement: %v", err)
	}

	return Img, nil
}

