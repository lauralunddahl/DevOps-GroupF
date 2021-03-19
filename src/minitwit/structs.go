package minitwit

import "time"

type User struct {
	Username string
	UserId   string
	Email    string
	PwHash   string
	Image    string
}

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
	PubDate   string
	Flagged   int
}

type Timeline struct {
	Username string
	UserId   int
	Email    string
	PwHash   string
	Image    string

	MessageId int
	AuthorId  int
	Text      string
	PubDate   time.Time
	Flagged   int
}