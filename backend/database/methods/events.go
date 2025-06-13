package sqlite

import "social-network/core/entities"

// EVENTS
func (d Database) GetAllEvents(CurrentUserID int) ([]entities.EventVM, error) {
	var events []entities.EventVM

	rows, err := d.db.Query(`
	SELECT DISTINCT
		e.ID,
		e.EventName,
		e.About,
		e.CreatedAt,
		g.GroupName,
		g.ID,
		COALESCE(
			(SELECT usr.Status 
			FROM EventUserRelations usr 
			WHERE e.ID = usr.EventID AND usr.UserID = ?),
			'pending'
		) AS CurrentStatus
		FROM Events AS e
		JOIN Groups g ON e.GroupID = g.ID
        JOIN GroupMemberRelations rel ON rel.GroupID = g.ID
        WHERE rel.UserID = ? AND (rel.GroupStatus = 2 OR rel.GroupStatus = 3)
		ORDER BY e.CreatedAt DESC;
		`, CurrentUserID, CurrentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e entities.EventVM
		var status string
		err = rows.Scan(&e.ID, &e.Name, &e.About, &e.Date, &e.GroupName, &e.GroupId, &status)
		if err != nil {
			return nil, err
		}
		if status == "pending" {
			e.Pending = 1
			e.Attending = 0
		} else if status == "notattending" {
			e.Pending = 0
			e.Attending = 0
		} else {
			e.Pending = 0
			e.Attending = 1
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (d Database) CreateEvent(CurrentUserID int, Event entities.EventDTO) ([]int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	stm, err := tx.Prepare("INSERT INTO Events(CreatorID, EventName, About, CreatedAt, GroupID) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = stm.Exec(CurrentUserID, Event.Name, Event.About, Event.Time, Event.GroupId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var users []int
	res, err := d.db.Query("SELECT UserID FROM GroupMemberRelations WHERE GroupID = ?", Event.GroupId)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var user int
		res.Scan(&user)
		users = append(users, user)
	}

	return users, nil
}
func (d Database) GetGroupEvents(CurrentUserID, GroupID int) ([]entities.EventVM, error) {
	var events []entities.EventVM
	query := `SELECT 
			e.ID,
			e.EventName,		
			e.About,
			e.CreatedAt,
			g.GroupName,
			g.ID,
			COALESCE(
				(SELECT usr.Status 
				FROM EventUserRelations usr 
				WHERE e.ID = usr.EventID AND usr.UserID = ?),
				'pending'
			) AS CurrentStatus
			FROM Events AS e
			JOIN Groups g ON e.GroupID = g.ID
			JOIN GroupMemberRelations rel ON rel.GroupID = g.ID
			WHERE 
			g.ID = ? AND rel.UserID = ?
			 AND (rel.GroupStatus = 2 OR rel.GroupStatus = 3)
			ORDER BY e.CreatedAt DESC;`

	rows, err := d.db.Query(query, CurrentUserID, GroupID, CurrentUserID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e entities.EventVM
		var status string
		err = rows.Scan(&e.ID, &e.Name, &e.About, &e.Date, &e.GroupName, &e.GroupId, &status)
		if err != nil {
			return nil, err
		}
		if status == "pending" {
			e.Pending = 1
		} else {
			e.Attending = 1
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (d Database) JoinEvent(CurrentUserID, EventID int) error {
	query := "INSERT INTO EventUserRelations(EventID, UserID, Status) VALUES (?,?,'attending')"

	if _, err := d.db.Exec(query, EventID, CurrentUserID); err != nil {
		return err
	}

	return nil
}

func (d Database) DenyEvent(CurrentUserID, EventID int) error {
	query := "INSERT INTO EventUserRelations(EventID, UserID, Status) VALUES (?,?,'notattending')"

	if _, err := d.db.Exec(query, EventID, CurrentUserID); err != nil {
		return err
	}

	return nil
}
