package structs

import "github.com/gorilla/websocket"

type User struct {
	ClientPhone string `json:"clientPhone"`
	ChatId      int    `json:"chatId"`
	Message     string `json:"message"`
	MessageId   int    `json:"messageId"`
	FileURL     string `json:"fileURL"`
	FileName    string `json:"fileName"`
	Messenger   string `json:"messenger"`
	conn        *websocket.Conn
}

type Messages struct {
	Msg chan *User
}
