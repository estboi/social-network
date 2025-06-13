package core

import (
	"errors"
	"social-network/core/entities"
)

func (c *Core) PostGetProc(PostId int) (entities.PostVM, error) {
	post, err := c.repo.PostRead(PostId)
	if err != nil {
		return entities.PostVM{}, err
	}
	return post, nil
}

func (c *Core) PostsGetGroupProc(GroupID, CurrentUserID int) ([]entities.PostVM, error) {
	posts, err := c.repo.GetGroupPosts(GroupID, CurrentUserID)
	if err != nil {
		return []entities.PostVM{}, err
	}
	return posts, nil
}

func (c *Core) PostsGetHomeProc(CurrentUserID int) ([]entities.PostVM, error) {
	posts, err := c.repo.GetHomePosts(CurrentUserID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *Core) PostsHandlerGetUser(UserID, currentUserId int) ([]entities.PostVM, error) {
	posts, err := c.repo.GetUserPosts(UserID, currentUserId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *Core) PostPostProc(CurrentUserID, GroupID int, dto entities.PostDTO) (int, error) {
	// Step 1: Validate
	if len(dto.Content) > 500 {
		return -1, errors.New("max content chars is 500")
	}
	var postId int
	var err error
	// Step 2: Check if its group related Post
	if GroupID != 0 {
		postId, err = c.repo.CreateGroupPost(CurrentUserID, GroupID, dto)
		if err != nil {
			return -1, err
		}
	} else {
		postId, err = c.repo.CreatePersonalPost(CurrentUserID, dto)
		if err != nil {
			return -1, err
		}
	}
	return postId, nil
}

// IMAGES PROCESSING
func (c *Core) PostsImagePostProc(PostId int, content []byte) error {
	if err := c.repo.RecordPostImage(PostId, content); err != nil {
		return err
	}
	return nil
}

func (c *Core) PostsImageGetProc(PostId int) ([]byte, error) {
	content, err := c.repo.ReadPostImage(PostId)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// ACTION PROCESSING
func (c *Core) PostsLikeProc(Votes entities.VotesData) (entities.VotesData, error) {

	// Step 1: Call the db method
	Votes, err := c.repo.LikePost(Votes)
	if err != nil {
		return Votes, err
	}

	return Votes, nil
}
