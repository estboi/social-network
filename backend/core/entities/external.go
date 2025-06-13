package entities

type Status struct {
	Message string `json:"message"`
}

type LoginDTO struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type RegisterDTO struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	NickName  *string `json:"nickName"`
	Email     string  `json:"email"`
	Pass      string  `json:"password"`
	Date      string  `json:"date"`
	Privacy   string  `json:"profilePrivacy"`
	About     string  `json:"about"`
}

type NavbarVM struct {
	Username string `json:"username"`
	UserID   int    `json:"userID"`
}

type PostDTO struct {
	Content string `json:"content"`
	Privacy string `json:"privacy"`
	Friends []int  `json:"friends"`
}

type CommentDTO struct {
	PostId    int    `json:"postId"`
	Content   string `json:"content"`
	UserId    int    `json:"creatorId"`
	CreatedAt string
	Votes     int
}

type UsersVM struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsFollower bool   `json:"isFollower"`
	IsFollow   bool   `json:"isFollow"`
	IsPending  bool   `json:"isPending"`
}

type UserFullVM struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Date      string `json:"date"`
	Avatar    []byte `json:"avatar"`
	About     string `json:"about"`
	Privacy   string `json:"privacy"`
	CanModify bool   `json:"canModify"`
}

type EventVM struct {
	ID        int    `json:id`
	Name      string `json:"name"`
	Date      string `json:"date"`
	About     string `json:"about"`
	Time      string `json:"time"`
	GroupName string `json:"groupName"`
	GroupId   int    `json:"groupId"`
	Attending int    `json:"attending"`
	Pending   int    `json:"pending"`
}

type EventDTO struct {
	Name    string `json:"name"`
	About   string `json:"about"`
	Time    string `json:"time"`
	GroupId int    `json:"groupId"`
}

type Message struct {
	MessageId   int    `json:"messageId"`
	ChatId      int    `json:"chatId"`
	ChatType    string `json:"type"`
	Content     string `json:"content"`
	SenderId    int    `json:"senderId"`
	Time        string `json:"time"`
	FirstName   string `json:"firstname"`
	SenderName  string `json:"sendername"`
	OtherUserId int    `json:"otheruser"`
}

type NotificationVM struct {
	UserId     int    `json:"userId"`
	SourceId   int    `json:"sourceId"`
	SourceType string `json:"sourceType"`
	NotifType  string `json:"type"`
	Content    string `jsong:"content"`
}

type PostVM struct {
	Id          int       `json:"postId"`
	CreatorVM   CreatorVM `json:"postCreator"`
	Content     string    `json:"postContent"`
	Time        string    `json:"postTime"`
	Votes       int       `json:"postVotes"`
	VotesTarget int       `json:"postVotesTarget"`
}

type CommentVM struct {
	Id          int       `json:"commentId"`
	CreatorVM   CreatorVM `json:"commentCreator"`
	Content     string    `json:"commentContent"`
	Time        string    `json:"commentTime"`
	Votes       int       `json:"commentVotes"`
	VotesTarget int       `json:"postVotesTarget"`
}

type CreatorVM struct {
	CreatorId   int    `json:"creatorId"`
	CreatorName string `json:"creatorName"`
}

type GroupVM struct {
	Id         int    `json:"id"`
	GroupName  string `json:"groupName"`
	GroupAbout string `json:"groupAbout"`
	IsMember   int    `json:"status"` // 0 - not member. 1 - pending. 2 - member. -1 - rejected. 3 - creator
}

type GroupDTO struct {
	GroupName  string `json:"groupName"`
	GroupAbout string `json:"groupAbout"`
}

type GroupsAction struct {
	UserID int `json:"userID"`
}
