package interfaces

import "social-network/core/entities"

// split to different interfaces in case of need
type Repository interface {
	//// Authentication

	// Returns user password hash if login correct,
	// or error if login is not correct or error with rCeading database.
	GetPasswordHash(userId int) (string, error)
	CreateUser(User entities.RegisterDTO) (int, error)
	LoginUser(entities.LoginDTO) (entities.LoginValues, error)

	//// Actions

	SubscribeOnUser(CurrentUserID, UserID int) (entities.Status, error)
	UnsubscribeOnUser(CurrentUserID, UserID int) (entities.Status, error)
	// RequestFriendship(CurrentUserID, UserID int) (entities.Status, error)
	// AcceptFriendship(CurrentUserID, UserID int) (entities.Status, error)

	//// Events

	//// Groups

	//// Images

	IsAccess2Image(CurrentUserID int, Src string, ID int) (bool, error)
	CreateImage(CurrentUserID int, Content []byte) error

	//// Navbar
	GetUserName(CurrentUserID int) (entities.NavbarVM, error)

	//// Notifications
	GetNotifications(CurrentUserID int) ([]entities.NotificationVM, error)
	NotificationRecord(VM entities.NotificationVM) error
	//// Posts

	IsAccess2ConcretePost(CurrentUserID, PostID int) (bool, error)
	IsAccess2GroupPosts(CurrentUserID, GroupID int) (bool, error)
	IsAccess2UserPosts(CurrentUserID, UserID int) (bool, error)

	// Posts interface
	PostsRepo
	// Comments interface
	CommentsRepo
	// Groups interface
	GroupsRepo
	// Events interface
	EventsRepo
	// User interface
	UserRepo
	// Chats interface
	ChatRepo
}

type ChatRepo interface {
	SaveMessage(msg entities.Message) error
	GetOneUserMessages(myUserId int, otherUserId int) []entities.Message
	GetOneGroupMessages(myUserId, groupID int) []entities.Message
	GetGroupChats(userId int) []entities.Message
	GetUserChats(userId int) []entities.Message
}

type UserRepo interface {
	GetAllUsers(CurrentUserID int) ([]entities.UsersVM, error)
	GetProfile(CurrentUserID, UserID int) (entities.UserFullVM, error)
	GetFollowed(CurrentUserID int) ([]entities.UsersVM, error)
	GetFollowers(CurrentUserID int) ([]entities.UsersVM, error)
	ModifyUser(int) (entities.Status, error)
	SubscribeOnUser(CurrentUserID, UserID int) (entities.Status, error)
	UnsubscribeOnUser(CurrentUserID, UserID int) (entities.Status, error)
	ImageUserRead(UserId int) ([]byte, error)
}

type PostsRepo interface {
	PostRead(PostId int) (entities.PostVM, error)
	GetGroupPosts(GroupID, CurrentUserID int) ([]entities.PostVM, error)
	GetHomePosts(CurrentUserID int) ([]entities.PostVM, error)
	GetUserPosts(UserID, CurrentUserID int) ([]entities.PostVM, error)
	// POSTS Images
	RecordPostImage(PostId int, Content []byte) error
	ReadPostImage(PostId int) ([]byte, error)

	CreatePersonalPost(CurrentUserID int, PostDTO entities.PostDTO) (int, error)
	CreateGroupPost(CurrentUserID, GroupID int, PostDTO entities.PostDTO) (int, error)
	// Posts actions
	LikePost(Votes entities.VotesData) (entities.VotesData, error)
}

type CommentsRepo interface {
	// Record
	CreateComment(CommentDTO entities.CommentDTO) (int, error)
	// Read
	GetComments(PostID int) ([]entities.CommentVM, error)
	// Images
	RecordCommentImage(commentId int, Content []byte) error
	ReadCommentImage(commentId int) ([]byte, error)
	// Actions
	LikeComment(Votes entities.VotesData) (entities.VotesData, error)
}

type GroupsRepo interface {
	// Record
	CreateGroup(CurrentUserID int, GroupDTO entities.GroupDTO) (int, error)
	// Read
	GetAllGroups(UserId int) ([]entities.GroupVM, error)
	GetCreatedGroups(CurrentUserID int) ([]entities.GroupVM, error)
	GetConnectedGroups(CurrentUserID int) ([]entities.GroupVM, error)
	GroupsProfileRead(UserId, GroupID int) (entities.GroupVM, error)
	GroupsGetNotMembers(groupId int) ([]entities.UsersVM, error)
	GroupsGetRequested(userId, groupId int) ([]entities.UsersVM, error)
	// Images
	GroupsRecordImage(groupId int, content []byte) error
	GroupsImageRead(groupId int) ([]byte, error)
	// Actions
	GroupRequest(CurrentUserID, GroupID int) (int, error)
	GroupInvite(CurrentUserID, GroupID int) error
	GroupAccept(CurrentUserID, GroupID int) error
	GroupDeny(CurrentUserID, GroupID int) error
}

type EventsRepo interface {
	// Read
	GetAllEvents(CurrentUserID int) ([]entities.EventVM, error)
	GetGroupEvents(CurrentUserID, GroupID int) ([]entities.EventVM, error)
	// Record
	CreateEvent(CurrentUserID int, EventDTO entities.EventDTO) ([]int, error)
	//ACTIONS
	JoinEvent(CurrentUserID, EventID int) error
	DenyEvent(CurrentUserID, EventID int) error
}
