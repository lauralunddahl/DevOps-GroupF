package dto

import database "db"

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
	PubDate   int
	Flagged   int
}

func AddMessage(author_id string, text string, pub_date int, flagged int) { //change pub_date to date at some point
	message := Message{AuthorId: author_id, Text: text, PubDate: pub_date, Flagged: flagged}
	result := database.DB.Create(&message)
	if result.Error != nil {
		print(result.Error)
	}
}
