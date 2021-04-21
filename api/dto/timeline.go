package dto

import (
	"strconv"
	"time"

	database "github.com/lauralunddahl/DevOps-GroupF/db"
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
var per_page = 30

func GetPrivateTimeline(user_id int) []Timeline {
	var timeline []Timeline
	res := db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("messages.flagged = 0 and (users.user_id = ? or users.user_id in (select whom_id from followers where who_id = ?))", strconv.Itoa(user_id), strconv.Itoa(user_id)).Order("messages.pub_date desc").Limit(per_page).Scan(&timeline)
	if res.Error != nil {
		log.Println("Private timeline")
		log.Error(res.Error)
	}
	return timeline
}

func GetPublicTimeline() []Timeline {
	var timeline []Timeline
	res := db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("messages.flagged = 0").Order("messages.pub_date desc").Limit(per_page).Scan(&timeline)
	if res.Error != nil {
		log.Println("Public timeline")
		log.Error(res.Error)
	}
	return timeline
}

func GetUserTimeline(profile_user_id int) []Timeline {
	var timeline []Timeline
	res := db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("users.user_id = ?", strconv.Itoa(profile_user_id)).Order("messages.pub_date desc").Limit(per_page).Scan(&timeline)
	if res.Error != nil {
		log.Println("User timeline")
		log.Error(res.Error)
	}
	return timeline
}
