package sqlite

import (
	"fmt"
	"log"
	"social-network/core/entities"
)

// CHATS
// TODO: Create chat function, isChat exists function for SaveMessage
func (d Database) SaveMessage(msg entities.Message) error {
	// prepare a SQL statement
	if msg.ChatType == "user" {
		stmt, err := d.db.Prepare("INSERT INTO UserChatMessages(SenderID, RecieverID, Content, Time) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(msg.SenderId, msg.ChatId, msg.Content, msg.Time)
		if err != nil {
			return err
		}
	} else if msg.ChatType == "group" {
		stmt, err := d.db.Prepare("INSERT INTO GroupChatMessages(SenderID, GroupID, Content, Time) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(msg.SenderId, msg.ChatId, msg.Content, msg.Time)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Database) GetUserChats(userId int) []entities.Message {
	messages := []entities.Message{}
	statement := `
		SELECT m.ID, m.Content, m.SenderID, m.RecieverID, m.Time, r.FirstName, s.FirstName
		FROM UserChatMessages m
		INNER JOIN Users r ON m.RecieverID = r.ID
		INNER JOIN Users s ON m.SenderID = s.ID
		WHERE m.SenderID = ? OR m.RecieverID = ?
		ORDER BY m.ID DESC
	`

	rows, err := d.db.Query(statement, userId, userId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	seen := make(map[int]bool) // To keep track of unique users
	for rows.Next() {
		newMessage := entities.Message{}
		err = rows.Scan(&newMessage.MessageId, &newMessage.Content, &newMessage.SenderId, &newMessage.ChatId, &newMessage.Time, &newMessage.FirstName, &newMessage.SenderName)
		if err != nil {
			fmt.Println(err)
		}
		newMessage.ChatType = "users"
		if newMessage.ChatId == userId {
			name, err := d.GetUserName(userId)
			if err != nil {
				log.Println(err)
			}
			newMessage.FirstName = name.Username
			newMessage.OtherUserId = newMessage.SenderId
		} else if newMessage.SenderId == userId {
			newMessage.OtherUserId = newMessage.ChatId
		}

		// Check if the user is already seen, if yes, update the message
		if _, ok := seen[newMessage.ChatId]; !ok && newMessage.SenderId == userId {
			messages = append(messages, newMessage)
			seen[newMessage.ChatId] = true
		}
	}

	return messages
}

func (d Database) GetGroupChats(userId int) []entities.Message {

	messages := []entities.Message{}
	statement := `
		SELECT m.ID, m.Content, m.SenderID, m.GroupID, m.Time, g.GroupName, s.FirstName
		FROM GroupChatMessages m
		INNER JOIN Groups g ON m.GroupID = g.ID
		INNER JOIN Users s ON m.SenderID = s.ID
		WHERE m.SenderID = ?
		ORDER BY m.ID DESC
	`

	rows, err := d.db.Query(statement, userId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	seen := make(map[int]bool) // To keep track of unique users
	for rows.Next() {
		newMessage := entities.Message{}
		err = rows.Scan(&newMessage.MessageId, &newMessage.Content, &newMessage.SenderId, &newMessage.ChatId, &newMessage.Time, &newMessage.FirstName, &newMessage.SenderName)
		if err != nil {
			fmt.Println(err)
		}
		newMessage.ChatType = "groups"
		// Check if the user is already seen, if yes, update the message
		if _, ok := seen[newMessage.ChatId]; !ok {
			messages = append(messages, newMessage)
			seen[newMessage.ChatId] = true
		}
	}

	return messages
}

func (d Database) GetOneUserMessages(myUserId int, otherUserId int) []entities.Message {
	singleMessages := []entities.Message{}
	statement := `
		SELECT m.Content, m.SenderID, m.RecieverID, m.Time, r.FirstName, s.FirstName
		FROM UserChatMessages m
		INNER JOIN Users r ON m.RecieverID = r.ID
		INNER JOIN Users s ON m.SenderID = s.ID
		WHERE (m.SenderID = ? AND m.RecieverID = ?) OR (m.SenderID = ? AND m.RecieverID = ?)
		ORDER BY m.ID DESC
	`

	rows, err := d.db.Query(statement, myUserId, otherUserId, otherUserId, myUserId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		newMessage := entities.Message{}
		err = rows.Scan(&newMessage.Content, &newMessage.SenderId, &newMessage.ChatId, &newMessage.Time, &newMessage.FirstName, &newMessage.SenderName)
		if err != nil {
			fmt.Println(err)
		}
		singleMessages = append(singleMessages, newMessage)
	}

	return singleMessages
}

func (d Database) GetOneGroupMessages(myUserId int, groupID int) []entities.Message {
	singleMessages := []entities.Message{}
	statement := `
		SELECT m.Content, m.SenderID, m.GroupID, m.Time, s.FirstName
		FROM GroupChatMessages m
		INNER JOIN Users s ON m.SenderID = s.ID
		WHERE m.GroupID = ?
		ORDER BY m.ID DESC
	`

	rows, err := d.db.Query(statement, groupID)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		newMessage := entities.Message{}
		err = rows.Scan(&newMessage.Content, &newMessage.SenderId, &newMessage.ChatId, &newMessage.Time, &newMessage.SenderName)
		if err != nil {
			fmt.Println(err)
		}
		singleMessages = append(singleMessages, newMessage)
	}

	return singleMessages
}
