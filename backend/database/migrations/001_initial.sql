-- 001_initial.sql

-- DO WE NEED TO: check txt column length 
--CHECK (LENGTH(NickName) <= 20)

-- DO WE NEED TO: check status or privacy string for example:
--CHECK (text_column IN ('Value1', 'Value2', 'Value3'))

-- DO WE NEED TO: on delete cascade somewhere?

-- Create Users table
CREATE TABLE IF NOT EXISTS Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    FirstName TEXT NOT NULL,
    LastName TEXT NOT NULL,
    NickName TEXT UNIQUE,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL,
    DateOfBirth TEXT NOT NULL,
    Img BLOB,
    About TEXT,
    Privacy TEXT  DEFAULT "Public"
);

-- Create UserSubscriptionRelations table
CREATE TABLE IF NOT EXISTS UserSubscriptionRelations (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER NOT NULL,
    FollowerID INTEGER NOT NULL,
    -- If private, we should wait for user response
    CurrentStatus TEXT NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (FollowerID) REFERENCES Users(ID)
);

-- Create Events table
CREATE TABLE IF NOT EXISTS Events (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    CreatorID INTEGER NOT NULL,
    EventName TEXT NOT NULL,
    About TEXT NOT NULL,
    CreatedAt TEXT NOT NULL,
    GroupID INTEGER NOT NULL,
    FOREIGN KEY (CreatorID) REFERENCES Users(ID)
    FOREIGN KEY (GroupID) REFERENCES Groups(ID)
);

-- Create EventUserRelations table
CREATE TABLE IF NOT EXISTS EventUserRelations (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    EventID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    Status TEXT DEFAULT "pending",
    FOREIGN KEY (EventID) REFERENCES Events(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);

-- Create Groups table
CREATE TABLE IF NOT EXISTS Groups (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    GroupName TEXT NOT NULL UNIQUE,
    GroupAbout TEXT NOT NULL,
    Img BLOB, 
    CreatorID INTEGER NOT NULL,
    FOREIGN KEY (CreatorID) REFERENCES Users(ID)
);

-- Create GroupMemberRelations table
CREATE TABLE IF NOT EXISTS GroupMemberRelations (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER NOT NULL,
    GroupID INTEGER NOT NULL,
    GroupStatus INTEGER DEFAULT 0, 
    FOREIGN KEY (UserID) REFERENCES Users(ID)
    FOREIGN KEY (GroupID) REFERENCES Groups(ID)
);

-- Create Posts table
CREATE TABLE IF NOT EXISTS Posts (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    CreatorID INTEGER NOT NULL,
    GroupID INTEGER, 
    Content TEXT NOT NULL,
    CreatedAt TEXT NOT NULL,
    Privacy TEXT NOT NULL,
    Votes INTEGER DEFAULT 0,
    Img BLOB,
    FOREIGN KEY (CreatorID) REFERENCES Users(ID)
    FOREIGN KEY (GroupID) REFERENCES Groups(ID)
);

CREATE TABLE IF NOT EXISTS PostsVotesRelation (
    PostID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    Target INTEGER DEFAULT 0,
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (PostID) REFERENCES Posts(ID)
);



-- Create Comments table
CREATE TABLE IF NOT EXISTS Comments (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    PostID INTEGER NOT NULL,
    CreatorID INTEGER NOT NULL,
    Content TEXT NOT NULL,
    CreatedAt TEXT NOT NULL,
    Votes INTEGER DEFAULT 0,
    Img BLOB,
    FOREIGN KEY (PostID) REFERENCES Posts(ID),
    FOREIGN KEY (CreatorID) REFERENCES Users(ID)
);

CREATE TABLE IF NOT EXISTS CommentsVotesRelation (
    CommentID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    Target INTEGER DEFAULT 0,
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (CommentID) REFERENCES Comments(ID)
);


-- Create ChatRelations table
CREATE TABLE IF NOT EXISTS GroupChatMessages (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    SenderID INTEGER NOT NULL,
    GroupID INTEGER NOT NULL,
    Content TEXT NOT NULL,
    Time TEXT NOT NULL,
    FOREIGN KEY (SenderID) REFERENCES Users(ID),
    FOREIGN KEY (GroupID) REFERENCES Groups(ID)
);

CREATE TABLE IF NOT EXISTS UserChatMessages (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    SenderID INTEGER NOT NULL,
    RecieverID INTEGER NOT NULL,
    Content TEXT NOT NULL,
    Time TEXT NOT NULL,
    FOREIGN KEY (SenderID) REFERENCES Users(ID),
    FOREIGN KEY (RecieverID) REFERENCES Users(ID)
);

-- Create Notifications table
CREATE TABLE IF NOT EXISTS Notifications (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    userID INTEGER NOT NULL,
    sourceID INTEGER NOT NULL,
    sourceType TEXT NOT NULL,
    type TEXT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (userID) REFERENCES Users(ID) 
);

CREATE TABLE IF NOT EXISTS PostAccess (
    userID INTEGER NOT NULL,
    postID INTEGER NOT NULL,
    FOREIGN KEY (userID) REFERENCES Users(ID)
    FOREIGN KEY (postID) REFERENCES Posts(ID)
)