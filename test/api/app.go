package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	wsConnect *websocket.Conn
}

type GlobalStorage struct {
	Users map[string]*User
	Mu    *sync.RWMutex
}

func GetGlobalStorage() *GlobalStorage {
	var gs GlobalStorage
	gs.initial()
	return &gs
}

func (gs *GlobalStorage) UserAdd(user *User) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()
	gs.Users[user.Phone] = user
}

func (gs *GlobalStorage) initial() {
	//TODO Обращаемся в БД за активными пользователями
	gs.Users = make(map[string]*User)
	gs.Mu = new(sync.RWMutex)
}

func worker() {
	for {
		conn, err := GetWSConnect()
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		msg := []byte("test message")
		postSocket1(conn, msg)
		break
	}
}

func main() {
	gs := GetGlobalStorage()

	conn, err := GetWSConnect()
	if err != nil {
		return
	}

	gs.Users["test"] = &User{
		Name:      "Alex",
		Phone:     "89068685856",
		wsConnect: conn,
	}

	for s, user := range gs.Users {
		fmt.Println(s, user)
		user.Phone = "627173676371237122731726"
	}
	fmt.Println(gs.Users["test"])
	worker()
}

var (
	upgrader = websocket.Upgrader{
		WriteBufferPool: &sync.Pool{},
		ReadBufferSize:  1024,
	}
)

type Message struct {
	ClientPhone string `json:"clientPhone"`
	ChatId      int    `json:"chatId"`
	Message     string `json:"message"`
	MessageId   int    `json:"messageId"`
	FileURL     string `json:"fileURL"`
	FileName    string `json:"fileName"`
	FileType    string `json:"fileType"`
}

func GetWSConnect() (*websocket.Conn, error) {
	var header = http.Header{}
	header.Add("Authorization", "token")
	conn, _, err := websocket.DefaultDialer.Dial("ws://172.24.18.44:8080/message/close", header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func postSocket1(conn *websocket.Conn, payload []byte) {
	err := conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		log.Println(err)
	}
	log.Println("Отправлено №2")
}
