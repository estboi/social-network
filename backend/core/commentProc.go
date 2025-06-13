package core

import (
	"social-network/core/entities"
)

func (c *Core) CommentPostProc(dto entities.CommentDTO) (int, error) {
	commentId, err := c.repo.CreateComment(dto)
	if err != nil {
		return -1, err
	}
	return commentId, nil
}
func (c *Core) CommentsGetProc(postId int) ([]entities.CommentVM, error) {
	// Step 1: get the comments from database
	comments, err := c.repo.GetComments(postId)
	if err != nil {
		return []entities.CommentVM{}, err
	}
	return comments, nil
}

// IMAGES PROCESSING
func (c *Core) CommentsImageCreateProc(commentId int, Content []byte) error {
	// Record image to database
	if err := c.repo.RecordCommentImage(commentId, Content); err != nil {
		return err
	}
	return nil
}

func (c *Core) CommentsImageGetProc(commentId int) ([]byte, error) {
	// Read image from database
	content, err := c.repo.ReadCommentImage(commentId)
	if err != nil {
		return nil, err
	}
	return content, err
}

// ACTION PROCESSING
func (c *Core) CommentsLikeProc(Votes entities.VotesData) (entities.VotesData, error) {
	var err error
	// Step 1: call the database method
	Votes, err = c.repo.LikeComment(Votes)
	if err != nil {
		return entities.VotesData{}, err
	}
	return Votes, nil
}
