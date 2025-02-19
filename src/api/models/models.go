package models

type Chatroom struct {
	ID              uint             `gorm:"primaryKey"`
	Sender_id       int
	Sender_name     *string
	Receiver_id     int
	Receiver_name   *string
	Messages_recent *string
	Messages_backup []MessageBackup `gorm:"type:jsonb" json:"messages_backup"`
}

type MessageBackup struct {
	IsSender bool   `json:"is_sender"`
	Message  string `json:"message"`
}

type Message struct {
	Sender_id   int
	Receiver_id int
	Message     string
}

type User struct {
	ID     uint `gorm:"primaryKey"`
	UserID int
	Name   string
}
