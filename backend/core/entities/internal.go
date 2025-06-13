package entities



type LoginValues struct {
	UserId       int
	HashPassword string
}

type VotesData struct {
	UserId int
	TargetId int
	TargetCreatorId int
	Votes  int
	Target int 
}

type EventDataAction struct{
	EventID int
}