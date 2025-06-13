package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/core/entities"
)

// GROUPS
func (d Database) GetConnectedGroups(currentUserID int) ([]entities.GroupVM, error) {
	rows, err := d.db.Query(`SELECT 
    g.ID,
    g.GroupName,
    g.GroupAbout,
    COALESCE(rel.GroupStatus, 0) AS Status
	FROM "Groups" AS g
	LEFT JOIN GroupMemberRelations AS rel ON rel.GroupID = g.ID 
	WHERE (rel.GroupStatus = 2 OR rel.GroupStatus = 3) AND rel.UserID = ?;
 	`, currentUserID)
	if err != nil {
		return nil, err
	}

	var groups []entities.GroupVM
	for rows.Next() {
		var g entities.GroupVM
		err := rows.Scan(&g.Id, &g.GroupName, &g.GroupAbout, &g.IsMember)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, g)
	}
	return groups, nil
}
func (d Database) GetCreatedGroups(CurrentUserID int) ([]entities.GroupVM, error) {
	query := `SELECT 
	g.ID, 
	g.GroupName,
	g.GroupAbout
	FROM Groups g WHERE g.CreatorID=?
	ORDER BY g.ID`

	rows, err := d.db.Query(query, CurrentUserID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	groups := make([]entities.GroupVM, 0)
	for rows.Next() {
		var group entities.GroupVM
		if err := rows.Scan(&group.Id, &group.GroupName, &group.GroupAbout); err != nil {
			return nil, fmt.Errorf("scanning row failed: %v", err)
		}
		group.IsMember = 3
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows contain error: %v", err)
	}

	return groups, nil
}
func (d Database) GetAllGroups(UserId int) ([]entities.GroupVM, error) {
	query := `SELECT 
    g.ID,
    g.GroupName,
    g.GroupAbout,
    COALESCE(rel.GroupStatus, 0) AS Status
	FROM "Groups" AS g
	LEFT JOIN GroupMemberRelations AS rel ON rel.GroupID = g.ID AND rel.UserID = ?
	WHERE (rel.UserID IS NULL OR rel.GroupStatus != -1);
 	`
	rows, err := d.db.Query(query, UserId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	groups := make([]entities.GroupVM, 0)
	for rows.Next() {
		var group entities.GroupVM
		if err := rows.Scan(&group.Id, &group.GroupName, &group.GroupAbout, &group.IsMember); err != nil {
			return nil, fmt.Errorf("scanning row failed: %v", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows contain error: %v", err)
	}

	return groups, nil
}

func (d Database) CreateGroup(CurrentUserID int, Group entities.GroupDTO) (int, error) {
	query := "INSERT INTO Groups(Groupname, GroupAbout, CreatorID) values(?,?,?)"
	relationQuery := "INSERT INTO GroupMemberRelations(UserID, GroupID, GroupStatus) values(?,?,3)"

	result, err := d.db.Exec(query, Group.GroupName, Group.GroupAbout, CurrentUserID)
	if err != nil {
		return -1, err
	}

	groupId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	if _, err = d.db.Exec(relationQuery, CurrentUserID, groupId); err != nil {
		return -1, err
	}
	return int(groupId), err
}

func (d Database) GroupsProfileRead(UserId, GroupID int) (entities.GroupVM, error) {
	query := `SELECT 
	g.ID, 
	g.GroupName, 
	g.GroupAbout,
	rel.GroupStatus
	FROM Groups AS g
	JOIN GroupMemberRelations AS rel ON g.ID = rel.GroupID 
	WHERE g.ID = ? AND rel.UserID = ? AND (rel.GroupStatus = 2 OR rel.GroupStatus = 3); `
	var group entities.GroupVM

	row := d.db.QueryRow(query, GroupID, UserId)
	err := row.Scan(&group.Id, &group.GroupName, &group.GroupAbout, &group.IsMember)

	if err != nil {
		if err == sql.ErrNoRows {
			// No result
			return entities.GroupVM{}, fmt.Errorf("no group found with ID: %d", GroupID)
		} else {
			// Different error
			return entities.GroupVM{}, err
		}
	}

	return group, nil
}

func (d Database) GroupsGetNotMembers(GroupID int) ([]entities.UsersVM, error) {
	query := `SELECT
    u.ID,
    COALESCE(u.NickName, u.FirstName || ' ' || u.LastName) AS UserName,
    COALESCE(rel.GroupStatus, 0) AS GroupStatus
	FROM Users AS u
	LEFT JOIN GroupMemberRelations AS rel ON u.ID = rel.UserID AND rel.GroupID = ?
	WHERE rel.GroupStatus = 4 OR rel.GroupStatus IS NULL;`

	var users []entities.UsersVM
	rows, err := d.db.Query(query, GroupID)
	if err != nil {
		return []entities.UsersVM{}, err
	}
	for rows.Next() {
		var u entities.UsersVM
		var pending int
		if err := rows.Scan(&u.ID, &u.Name, &pending); err != nil {
			log.Fatal(err)
		}
		if pending == 4 {
			u.IsPending = true
		}
		users = append(users, u)
	}
	return users, nil
}

func (d Database) GroupsGetRequested(UserId, GroupId int) ([]entities.UsersVM, error) {
	query := `SELECT 
	u.ID,
	COALESCE(u.NickName, u.FirstName || ' ' || u.LastName) AS UserName
	FROM GroupMemberRelations as rel
	JOIN Users AS u ON rel.UserID = u.ID
	WHERE rel.GroupStatus = 1 AND rel.GroupID = ?;
	`
	var users []entities.UsersVM
	rows, err := d.db.Query(query, GroupId)
	if err != nil {
		return []entities.UsersVM{}, err
	}
	for rows.Next() {
		var u entities.UsersVM
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	return users, nil
}

// ACTIONS
func (d Database) GroupRequest(CurrentUserID, GroupID int) (int, error) {
	stmt, err := d.db.Prepare("INSERT INTO GroupMemberRelations(UserID, GroupID, GroupStatus) VALUES(?, ?, ?)")
	if err != nil {
		return -1, err
	}

	if _, err = stmt.Exec(CurrentUserID, GroupID, 1); err != nil {
		return -1, err
	}

	var creatorID int
	err = d.db.QueryRow("SELECT CreatorID FROM Groups WHERE ID = ?", GroupID).Scan(&creatorID)
	if err != nil {
		return -1, err
	}

	return creatorID, nil
}

func (d Database) GroupInvite(InvitedUserId, GroupID int) error {
	query := "INSERT INTO GroupMemberRelations(UserID, GroupID, GroupStatus) VALUES (?,?,4)"
	if _, err := d.db.Exec(query, InvitedUserId, GroupID); err != nil {
		return err
	}
	return nil
}

func (d Database) GroupAccept(RequestingUserId, GroupID int) error {
	query := `UPDATE GroupMemberRelations SET GroupStatus = 2 WHERE UserID = ? AND GroupID = ?`
	if _, err := d.db.Exec(query, RequestingUserId, GroupID); err != nil {
		return err
	}
	return nil
}
func (d Database) GroupDeny(RequestingUserId, GroupID int) error {
	query := `UPDATE GroupMemberRelations SET GroupStatus = -1 WHERE UserID = ? AND GroupID = ?`
	if _, err := d.db.Exec(query, RequestingUserId, GroupID); err != nil {
		return err
	}
	return nil
}

// IMAGES

func (d Database) GroupsRecordImage(groupId int, content []byte) error {
	query := "UPDATE Groups SET Img = ? WHERE ID = ?"
	if _, err := d.db.Exec(query, content, groupId); err != nil {
		return err
	}
	return nil
}

func (d Database) GroupsImageRead(groupId int) ([]byte, error) {
	query := "SELECT Img FROM Groups WHERE ID = ?"
	var content []byte

	if err := d.db.QueryRow(query, groupId).Scan(&content); err != nil {
		return nil, err
	}

	return content, nil
}
