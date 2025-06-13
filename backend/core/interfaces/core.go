package interfaces

import (
	"social-network/core/entities"
)

type Core interface {
	//// Authentication Handlers

	// Anonymous user sends credentials to be logged in
	// Success in case password match with hash.
	LoginProcess(entities.LoginDTO) (int, error)

	// Anonymous user sends his data to be registered.
	// Success in case data is valid.
	RegisterProc(entities.RegisterDTO) (int, error)

	//// Action Handlers
	//
	SubscribeProcessing(CurrentUserID, UserID int) (entities.Status, error)

	ModifyUserProcessing(CurrentUserID int) (entities.Status, error)

	//
	UnsubscribeProcessing(CurrentUserID, UserID int) (entities.Status, error)

	// Current user accepts group invitation.

	// Current user (group owner) accepts request for joining his group.

	//// Chats Handlers
	SendMessage(msg entities.Message) error

	// Current user trying to receive all chats with both users and groups he chatting in.
	// Core returns all user chats.
	GetChatMessages(CurrentUserID, chatID int) ([]entities.Message, error)

	// Current user trying to get chat history within group chat by group id.
	// If user is a group member, core returns history of messages.
	ChatGetAllMessagesProc(CurrentUserID int) ([]entities.Message, error)

	ChatGetGroupMessagesProc(CurrentUserID, groupID int) ([]entities.Message, error)

	// Current user trying to receive img from specific source table by id.
	// Core checks is source exists and is user have access to it.
	// In case of correct source read and return img by id.
	ImageUserProc(UserId int) ([]byte, error)

	// Current user trying to save img to specific source source table by id.
	// Core checks is source exists and is user have access to it.
	// In case of correct source write img to it by id.
	ImageUserCreateProc(CurrentUserID int, Content []byte) error

	//// Navbar Handler
	// Just returning username by user ID
	NavbarProc(UserID int) (entities.NavbarVM, error)

	//// Notifications Handler
	// Just return all notifications for current user.
	NotificationsHandlerGet(CurrentUserID int) ([]entities.NotificationVM, error)
	NotificationRecord(notif entities.NotificationVM) error

	//// Users Handlers
	// Just returning all users on a platform.
	GetAllUsersProcessing(userId, From, To int) ([]entities.UsersVM, error)

	// Current user request users he follows.
	// Core processing it without any checks.
	// Returning all followed users.
	GetFollowedUsersProcessing(CurrentUserID int) ([]entities.UsersVM, error)

	// Current user request own followers.
	// Core processing it without any checks.
	// Returning all followers users.
	GetFollowersProcessing(CurrentUserID, From, To int) ([]entities.UsersVM, error)

	// Current user request another user profile data.
	// Core processing it without any checks.
	// Returning data as the database provides.
	GetUserProfileProcessing(CurrentUserID, UserID int) (entities.UserFullVM, error)

	// Posts interface
	Posts
	// Comments interface
	Comments
	// Groups interface
	Groups
	// Events interface
	Events
}

type Posts interface {
	// Main proccessing

	// Get one post
	PostGetProc(PostID int) (entities.PostVM, error)
	PostPostProc(CurrentUserID, GroupID int, PostDTO entities.PostDTO) (int, error)
	PostsGetHomeProc(CurrentUserID int) ([]entities.PostVM, error)
	PostsGetGroupProc(GroupID, CurrentUserID int) ([]entities.PostVM, error)
	PostsHandlerGetUser(UserID, CurrentUserID int) ([]entities.PostVM, error)

	// Images processing
	PostsImagePostProc(PostID int, Content []byte) error
	PostsImageGetProc(PostId int) ([]byte, error)

	// Action processing
	PostsLikeProc(Votes entities.VotesData) (entities.VotesData, error)
}

type Comments interface {
	// Main processing
	CommentsGetProc(PostID int) ([]entities.CommentVM, error)
	CommentPostProc(CommentDTO entities.CommentDTO) (int, error)

	// Image processing
	CommentsImageCreateProc(commentId int, Content []byte) error
	CommentsImageGetProc(commentId int) ([]byte, error)

	// Action processing
	CommentsLikeProc(Votes entities.VotesData) (entities.VotesData, error)
}

type Groups interface {
	// Get processing
	GroupsAllProc(userId int) ([]entities.GroupVM, error)
	GroupsHandlerGetConnected(CurrentUserID int) ([]entities.GroupVM, error)
	GroupsHandlerGetCreated(CurrentUserID int) ([]entities.GroupVM, error)
	GroupsProfileProc(UserId, groupID int) (entities.GroupVM, error)

	// Options processing
	GroupsGetNotMembersProc(groupId int) ([]entities.UsersVM, error)
	GroupsGetRequestedProc(userId, groupId int) ([]entities.UsersVM, error)
	// Main processing
	GroupsCreateProc(CurrentUserID int, GroupDTO entities.GroupDTO) (int, error)
	// Action Processing
	GroupsRequestProc(CurrentUserID, GroupID int) (int, error)
	GroupsInviteProc(InvitedUserId, GroupID int) error
	GroupsAcceptProc(CurrentUserID, GroupID int) error
	GroupsDenyProc(RequestingUserId, GroupID int) error

	// Image processing
	GroupsImageCreateProc(grouId int, content []byte) error
	GroupsImageProc(grouId int) ([]byte, error)
}

type Events interface {
	// Get Processing
	EventGetAllProc(CurrentUserID int) ([]entities.EventVM, error)
	EventGetGroupProc(CurrentUserID, GroupID int) ([]entities.EventVM, error)
	// Create processing
	EventCreateProc(CurrentUserID int, EventDTO entities.EventDTO) ([]int, error)
	// Action processing
	EventAcceptProc(CurrentUserID, EventID int) error
	EventDenyProc(CurrentUserID, EventID int) error
}
