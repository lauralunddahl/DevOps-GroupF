package dto

import (
	database "db"
	"strconv"

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
	PubDate   int
	Flagged   int
}

var db *gorm.DB = database.DB
var per_page = 30

// func initialMigration() {
// 	db.AutoMigrate(&Timeline{})
// }

func GetPrivateTimeline(user_id int) []Timeline {
	var timeline []Timeline
	db.Table("message").Select("user.*, message.*").Joins("join user on message.author_id = user.user_id").Where("message.flagged = 0 and (user.user_id = ? or user.user_id in (select whom_id from follower where who_id = ?))", strconv.Itoa(user_id), strconv.Itoa(user_id)).Order("message.pub_date desc").Limit(per_page).Scan(&timeline)
	return timeline
}
