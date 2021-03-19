package api

import "time"

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

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
	PubDate   time.Time
	Flagged   int
}

type Latest struct {
	Latest int `json:"latest"`
}

type Response struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"error_msg"`
}

type ApiMessage struct {
	Content string    `json:"content"`
	PubDate time.Time `json:"pub_date"`
	User    string    `json:"user"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email" `
	Password string `json:"pwd"`
}

type Followers struct {
	Username string `json:"username"`
}
type FollowUser struct {
	Follow   string `json:"follow"`
	Unfollow string `json:"unfollow"`
}