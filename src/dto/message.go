package dto

import (
	"time"

	database "github.com/lauralunddahl/DevOps-GroupF/src/db"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	MessageId int
	AuthorId  string
	Text      string
	PubDate   time.Time
	Flagged   int
}

func AddMessage(author_id string, text string, pub_date time.Time, flagged int) { //change pub_date to date at some point
	message := Message{AuthorId: author_id, Text: text, PubDate: pub_date, Flagged: flagged}
	result := database.DB.Create(&message)
	if result.Error != nil {
		log.Println("AddMessage")
		log.Error(result.Error)
	}
}


func GetTotalNumberOfMessages() int64 {
	var result int64
	database.DB.Table("messages").Count(&result)
	return result
}
