package dto

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
	PubDate   string
	Flagged   int
}

func (Message) TableName() string {
	return "message"
}

// func initialMigration() {
// 	db.AutoMigrate(&Message{})
// }
