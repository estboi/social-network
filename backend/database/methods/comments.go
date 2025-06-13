package sqlite

import (
	"errors"
	"log"
	"social-network/core/entities"
	"time"
)

// COMMENTS
func (d Database) GetComments(PostID int) ([]entities.CommentVM, error) {
	var comments []entities.CommentVM
	query := `
	SELECT 
		C.ID,
		C.CreatorID,
		COALESCE(Users.NickName, Users.FirstName || ' ' || Users.LastName) AS UserName,
		C.Content,
		C.CreatedAt,
		C.Votes,
		COALESCE(CommentsVotesRelation.Target, 0) AS Target
	FROM 
		Comments AS C
	JOIN 
		Users ON C.CreatorID = Users.ID
	LEFT JOIN 
		CommentsVotesRelation ON C.ID = CommentsVotesRelation.CommentID
	WHERE C.PostID = ?
	ORDER BY
		C.CreatedAt DESC;`
	rows, err := d.db.Query(query, PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment entities.CommentVM
		var creator entities.CreatorVM
		err := rows.Scan(
			&comment.Id,
			&creator.CreatorId,
			&creator.CreatorName,
			&comment.Content,
			&comment.Time,
			&comment.Votes,
			&comment.VotesTarget)
		if err != nil {
			return nil, err
		}
		comment.CreatorVM = creator

		comments = append(comments, comment)
	}
	return comments, nil
}
func (d Database) CreateComment(Comment entities.CommentDTO) (int, error) {
	// Define the current timestamp.
	Comment.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	// Prepare the SQL statement for inserting a new comment.
	stmt, err := d.db.Prepare("INSERT INTO Comments (PostID, CreatorID, Content, CreatedAt) VALUES (?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// Execute the SQL statement with the provided parameters.
	result, err := stmt.Exec(Comment.PostId, Comment.UserId, Comment.Content, Comment.CreatedAt)
	if err != nil {
		return -1, err
	}
	commentId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(commentId), nil
}
func (d Database) LikeComment(Votes entities.VotesData) (entities.VotesData, error) {
	// Check if the post exists
	var commentExists int
	err := d.db.QueryRow("SELECT COUNT(*) FROM Comments WHERE ID = ?", Votes.TargetId).Scan(&commentExists)
	if err != nil {
		return entities.VotesData{}, err
	}
	if commentExists != 1 {
		return entities.VotesData{}, errors.New("comment does not exist")
	}

	// Increase the votes count of the post
	if _, err = d.db.Exec(`UPDATE Comments SET Votes = ? WHERE ID = ?`, Votes.Votes, Votes.TargetId); err != nil {
		return entities.VotesData{}, err
	}

	var relationExists int
	// Store the previous interaction with this post to prevent double action
	if err = d.db.QueryRow(`SELECT COUNT(*) FROM CommentsVotesRelation 
	WHERE CommentID = ? AND UserID = ? `, Votes.TargetId, Votes.UserId).Scan(&relationExists); err != nil {
		return entities.VotesData{}, err
	}

	if relationExists == 1 {
		var target int
		query := "SELECT Target FROM CommentsVotesRelation WHERE CommentID = ? AND UserID = ?"
		if err := d.db.QueryRow(query, Votes.TargetId, Votes.UserId).Scan(&target); err != nil {
			return entities.VotesData{}, err
		}
		if target != Votes.Target && target != 0 {
			Votes.Target = 0
		}
		query = "UPDATE CommentsVotesRelation SET Target = ? WHERE CommentID = ? AND UserID = ?"
		if _, err := d.db.Exec(query, Votes.Target, Votes.TargetId, Votes.UserId); err != nil {
			return entities.VotesData{}, err
		}
		return Votes, err
	}
	if _, err = d.db.Exec(`INSERT INTO CommentsVotesRelation (CommentID, UserID, Target) 
		VALUES (?,?,?)`, Votes.TargetId, Votes.UserId, Votes.Target); err != nil {
		log.Println(err)
		return Votes, err
	}
	return Votes, err
}

// IMAGES
func (d Database) RecordCommentImage(CommentId int, Content []byte) error {
	query := "UPDATE Comments SET Img = ? WHERE ID = ?"
	if _, err := d.db.Exec(query, Content, CommentId); err != nil {
		return err
	}
	return nil
}

func (d Database) ReadCommentImage(commentId int) ([]byte, error) {
	var VM []byte
	query := "SELECT Img FROM Comments WHERE ID = ?"
	if err := d.db.QueryRow(query, commentId).Scan(&VM); err != nil {
		return nil, err
	}
	return VM, nil
}
