package sqlite

import "social-network/core/entities"

// NOTIFICATIONS
func (d Database) GetNotifications(CurrentUserID int) ([]entities.NotificationVM, error) {
	var notificationList []entities.NotificationVM
	query := `
	SELECT 
	userID,
	sourceID,
	sourceType,
	type,
	content
	FROM "Notifications"
	WHERE userID = ?
	ORDER BY ID DESC;
	`

	rows, err := d.db.Query(query, CurrentUserID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var notification entities.NotificationVM
		err := rows.Scan(&notification.UserId, &notification.SourceId, &notification.SourceType, &notification.NotifType, &notification.Content)
		if err != nil {
			return nil, err
		}
		notificationList = append(notificationList, notification)
	}
	return notificationList, nil
}

func (d Database) NotificationRecord(VM entities.NotificationVM) error {
	query := `INSERT INTO Notifications 
	(userID, sourceID, sourceType, type, content) 
	VALUES (?,?,?,?,?)`

	_, err := d.db.Exec(query, VM.UserId, VM.SourceId, VM.SourceType, VM.NotifType, VM.Content)
	if err != nil {
		return err
	}

	return nil
}
