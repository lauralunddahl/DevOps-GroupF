package dto

import (
	database "db"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Timeline struct {
	Username string
	UserId   int
	Email    string
	PwHash   string

	MessageId int
	AuthorId  int
	Text      string
	PubDate   time.Time
	Flagged   int
}

var db *gorm.DB = database.DB
var per_page = 30

func GetPrivateTimeline(user_id int) []Timeline {
	var timeline []Timeline
	db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("messages.flagged = 0 and (users.user_id = ? or users.user_id in (select whom_id from followers where who_id = ?))", strconv.Itoa(user_id), strconv.Itoa(user_id)).Order("messages.pub_date desc").Limit(per_page).Scan(&timeline)
	return timeline
}

func GetPublicTimeline() []Timeline {
	var timeline []Timeline
	db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("messages.flagged = 0").Order("messages.pub_date desc").Limit(per_page).Scan(&timeline)
	return timeline
}

func GetUserTimeline(profile_user_id int) []Timeline {
	var timeline []Timeline
	db.Table("messages").Select("users.*, messages.*").Joins("join users on messages.author_id = users.user_id").Where("users.user_id = ?", strconv.Itoa(profile_user_id)).Order("messages.pub_date desc").Limit(per_page).Scan(&timeline)
	return timeline
}
