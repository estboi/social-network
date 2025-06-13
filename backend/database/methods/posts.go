package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social-network/core/entities"
	"time"
)

// Posts

func (d Database) PostRead(PostID int) (entities.PostVM, error) {
	query := `
	SELECT 
    	Posts.ID,
    	Posts.CreatorId,
    	COALESCE(Users.NickName, Users.FirstName || ' ' || Users.LastName) AS UserName,
    	Posts.Content,
    	Posts.CreatedAt,
    	Posts.Votes,
    	COALESCE(PostsVotesRelation.Target, 0) AS Target
	FROM Posts 
	JOIN Users ON Posts.CreatorID = Users.ID
	LEFT JOIN 
    PostsVotesRelation ON Posts.ID = PostsVotesRelation.PostID
	WHERE Posts.ID = ?
	`
	var post entities.PostVM
	if err := d.db.QueryRow(query, PostID).Scan(&post.Id,
		&post.CreatorVM.CreatorId,
		&post.CreatorVM.CreatorName,
		&post.Content,
		&post.Time,
		&post.Votes,
		&post.VotesTarget); err != nil {
		return entities.PostVM{}, err
	}

	return post, nil
}

func (d Database) GetHomePosts(CurrentUserID int) ([]entities.PostVM, error) {
	var posts []entities.PostVM

	query := `
	SELECT DISTINCT
    Posts.ID,
    Posts.CreatorId,
    COALESCE(Users.NickName, Users.FirstName || ' ' || Users.LastName) AS UserName,
    Posts.Content,
    Posts.CreatedAt,
    Posts.Votes,
    COALESCE(rel.Target, 0) AS Target
FROM 
    Posts 
JOIN 
    Users ON Posts.CreatorID = Users.ID
LEFT JOIN 
    PostsVotesRelation AS rel ON Posts.ID = rel.PostID AND rel.UserID = ?
WHERE 
    Posts.Privacy = 'Public' 
AND
    Posts.GroupID IS NULL
ORDER BY 
    Posts.CreatedAt DESC;
	`
	rows, err := d.db.Query(query, CurrentUserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post := entities.PostVM{}
		err := rows.Scan(
			&post.Id,
			&post.CreatorVM.CreatorId,
			&post.CreatorVM.CreatorName,
			&post.Content,
			&post.Time,
			&post.Votes,
			&post.VotesTarget,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (d Database) GetGroupPosts(GroupID, CurrentUserID int) ([]entities.PostVM, error) {
	rows, err := d.db.Query(`
	SELECT p.ID, 
	u.ID, 
	COALESCE(u.NickName, u.FirstName || ' ' || u.LastName) AS UserName,
	p.Content, 
	p.CreatedAt,
	p.Votes,
    COALESCE(PostsVotesRelation.Target, 0) AS Target
	FROM Posts P 
	LEFT JOIN PostsVotesRelation ON p.ID = PostsVotesRelation.PostID AND PostsVotesRelation.UserID = ? 
	JOIN Users u on p.CreatorID = u.ID 
	WHERE p.GroupID = ?
	ORDER BY p.CreatedAt DESC;
	`, CurrentUserID, GroupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entities.PostVM

	for rows.Next() {
		var post entities.PostVM
		var time string

		err := rows.Scan(
			&post.Id,
			&post.CreatorVM.CreatorId,
			&post.CreatorVM.CreatorName,
			&post.Content,
			&post.Time,
			&post.Votes,
			&post.VotesTarget)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		post.Time = time
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (d Database) GetUserPosts(UserID, CurrentUserID int) ([]entities.PostVM, error) {
	var posts []entities.PostVM

	rows, err := d.db.Query(`
	SELECT
	P.ID, 
	U.ID,
	COALESCE(U.NickName, U.FirstName || ' ' || U.LastName) AS UserName,
	P.Content,
	P.CreatedAt, 
	P.Votes,
	COALESCE(PostsVotesRelation.Target, 0) AS Target,
	P.Privacy
	FROM Posts P 
	LEFT JOIN PostsVotesRelation ON p.ID = PostsVotesRelation.PostID AND PostsVotesRelation.UserID = ? 
	LEFT JOIN Users U ON P.CreatorID = U.ID 
	WHERE P.CreatorID = ? AND P.GroupID IS NULL
	ORDER BY P.CreatedAt DESC;`, CurrentUserID, UserID)
	if err != nil {
		return nil, errors.New("failed to execute database query")
	}
	defer rows.Close()

	for rows.Next() {
		var post entities.PostVM
		var privacy string
		var timeStr string
		err = rows.Scan(&post.Id, &post.CreatorVM.CreatorId, &post.CreatorVM.CreatorName, &post.Content, &timeStr, &post.Votes, &post.VotesTarget, &privacy)
		if err != nil {
			return nil, errors.New("failed to scan row")
		}
		// Validate the access to user
		if UserID != CurrentUserID && privacy != "Public" {
			hasAccess, err := d.db.Query("SELECT 1 FROM PostAccess WHERE postID = ? AND userID = ?", post.Id, CurrentUserID, UserID)
			if err != nil {
				return nil, err
			}
			if !hasAccess.Next() {
				continue
			}
		}

		post.Time = timeStr
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("failed to iterate over rows")
	}

	return posts, nil
}

func (d Database) LikePost(Votes entities.VotesData) (entities.VotesData, error) {
	// Check if the post exists
	var postExists int
	err := d.db.QueryRow("SELECT COUNT(*) FROM Posts WHERE ID = ?", Votes.TargetId).Scan(&postExists)
	if err != nil {
		return entities.VotesData{}, err
	}
	if postExists != 1 {
		return entities.VotesData{}, errors.New("post does not exist")
	}

	// Increase the votes count of the post
	if _, err = d.db.Exec(`UPDATE Posts SET Votes = ? WHERE ID = ?`, Votes.Votes, Votes.TargetId); err != nil {
		return entities.VotesData{}, err
	}

	var relationExists int

	if err = d.db.QueryRow(`SELECT COUNT(*) FROM PostsVotesRelation 
	WHERE PostID = ? AND UserID = ? `, Votes.TargetId, Votes.UserId).Scan(&relationExists); err != nil {
		return entities.VotesData{}, err
	}

	if relationExists == 1 {
		var target int
		query := "SELECT Target FROM PostsVotesRelation WHERE PostID = ? AND UserID = ?"
		if err := d.db.QueryRow(query, Votes.TargetId, Votes.UserId).Scan(&target); err != nil {
			return entities.VotesData{}, err
		}
		if target != Votes.Target && target != 0 {
			Votes.Target = 0
		}
		query = "UPDATE PostsVotesRelation SET Target = ? WHERE PostID = ? AND UserID = ?"
		if _, err := d.db.Exec(query, Votes.Target, Votes.TargetId, Votes.UserId); err != nil {
			return entities.VotesData{}, err
		}

		d.db.QueryRow("SELECT CreatorID FROM Posts WHERE ID =?", Votes.TargetId).Scan(&Votes.TargetCreatorId)
		return Votes, err
	}
	if _, err = d.db.Exec(`INSERT INTO PostsVotesRelation (PostID, UserID, Target) 
		VALUES (?,?,?)`, Votes.TargetId, Votes.UserId, Votes.Target); err != nil {
		return Votes, err
	}
	return Votes, err
}

// Create posting
func (d Database) CreatePersonalPost(CurrentUserID int, Post entities.PostDTO) (int, error) {
	// Define the current timestamp.
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Create the personal post.
	insertQuery := "INSERT INTO Posts (CreatorID, Content, CreatedAt, Privacy) VALUES (?, ?, ?, ?);"
	result, err := d.db.Exec(insertQuery, CurrentUserID, Post.Content, currentTime, Post.Privacy)
	if err != nil {
		return -1, err
	}

	// Get the Id of created post
	postId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	// If Privacy == Followers
	if Post.Privacy == "Followers" && len(Post.Friends) != 0 {
		query := "INSERT INTO PostAccess(userID, postID) VALUES (?,?)"
		for _, friendID := range Post.Friends {
			if _, err := d.db.Exec(query, friendID, postId); err != nil {
				return -1, err
			}
		}
	}

	return int(postId), nil
}

func (d Database) CreateGroupPost(CurrentUserID, GroupID int, Post entities.PostDTO) (int, error) {
	// Insert the post
	result, err := d.db.Exec("INSERT INTO Posts (CreatorID, GroupID, Content, CreatedAt, Privacy, Votes) VALUES (?, ?, ?, ?, ?, ?)",
		CurrentUserID, GroupID, Post.Content, time.Now().Format(time.RFC3339), Post.Privacy, 0)
	if err != nil {
		return -1, err
	}
	// Get the Id of created post
	postId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(postId), err
}

// Actions
func (d Database) IsAccess2ConcretePost(CurrentUserID, PostID int) (bool, error) {
	var creatorID, groupID int
	row := d.db.QueryRow("SELECT CreatorID, GroupID FROM Posts WHERE ID = ?", PostID)
	err := row.Scan(&creatorID, &groupID)
	if err != nil {
		return false, fmt.Errorf("failed to execute SELECT query: %w", err)
	}
	// Current user created the post
	if CurrentUserID == creatorID {
		return true, nil
	}
	// Current user is in the group, that this post belongs to
	var userID int
	err = d.db.QueryRow("SELECT UserID FROM GroupMemberRelations WHERE UserID = ? AND GroupID = ?", CurrentUserID, groupID).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to execute SELECT query: %w", err)
	} else if err == nil && userID == CurrentUserID {
		return true, nil
	}
	return false, nil
}

func (d Database) IsAccess2GroupPosts(CurrentUserID, GroupID int) (bool, error) {

	query := `SELECT COUNT(*) FROM GroupMemberRelations WHERE UserID = ? AND GroupID = ?`

	var numRecords int

	err := d.db.QueryRow(query, CurrentUserID, GroupID).Scan(&numRecords)

	if err != nil {
		fmt.Println("Failed to execute query.")
		return false, err
	}

	if numRecords > 0 {
		return true, nil
	}

	return false, nil
}

func (d Database) IsAccess2UserPosts(CurrentUserID, UserID int) (bool, error) {
	// Prepare a SQL query
	stmt, err := d.db.Prepare(`SELECT COUNT(*) FROM UserSubscriptionRelations WHERE UserID = ? AND FollowerID = ?`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the query
	var count int
	err = stmt.QueryRow(UserID, CurrentUserID).Scan(&count)
	if err != nil {
		return false, err
	}

	// Check if there is a subscription relation between two users
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
