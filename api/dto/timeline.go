package dto

import (
	"strconv"
	"time"
	database "github.com/lauralunddahl/DevOps-GroupF/api/db"
	log "github.com/sirupsen/logrus"
)

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

var db = database.DB
var perPage = 30

func GetPrivateTimeline(userId int) []Timeline {
	var timeline []Timeline
	res := db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("messages.flagged = 0 and (users.user_id = ? or users.user_id in (select whom_id from followers where who_id = ?))", strconv.Itoa(userId), strconv.Itoa(userId)).Order("messages.pub_date desc").Limit(perPage).Scan(&timeline)
	if res.Error != nil {
		log.Println("Private timeline")
		log.Error(res.Error)
	}
	return timeline
}

func GetPublicTimeline() []Timeline {
	var timeline []Timeline
	res := db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("messages.flagged = 0").Order("messages.pub_date desc").Limit(perPage).Scan(&timeline)
	if res.Error != nil {
		log.Println("Public timeline")
		log.Error(res.Error)
	}
	return timeline
}

func GetUserTimeline(profileUserId int) []Timeline {
	var timeline []Timeline
	res := db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("users.user_id = ?", strconv.Itoa(profileUserId)).Order("messages.pub_date desc").Limit(perPage).Scan(&timeline)
	if res.Error != nil {
		log.Println("User timeline")
		log.Error(res.Error)
	}
	return timeline
}
