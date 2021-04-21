package dto

import (
	"time"

	database "github.com/lauralunddahl/DevOps-GroupF/api/db"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	MessageId int
	AuthorId  string
	Text      string
	PubDate   time.Time
	Flagged   int
}

func AddMessage(authorId string, text string, pubDate time.Time, flagged int) {
	message := Message{AuthorId: authorId, Text: text, PubDate: pubDate, Flagged: flagged}
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
