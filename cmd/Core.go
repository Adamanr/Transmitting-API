package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"test/config"
	"test/mainworker/structs"
	"test/socket"
	"time"
)

type GlobalStorage struct {
	Users map[string]*structs.User
	Mu    *sync.RWMutex
}

func GetGlobalStorage() *GlobalStorage {
	var gs GlobalStorage
	gs.initial()
	return &gs
}

func (gs *GlobalStorage) initial() {
	gs.Users = make(map[string]*structs.User)
	gs.Mu = new(sync.RWMutex)
}

// UserAdd Добавляем пользователя в стуктуру
func (gs *GlobalStorage) UserAdd(user *structs.User) []byte {
	gs.Mu.Lock()
	gs.Users[user.ClientPhone] = user
	gs.Mu.Unlock()

	if users, err := json.Marshal(user); err != nil {
		log.Println(err)
		return nil
	} else {
		return users
	}
}

// GetUsers Получает данные от клиента и отправляет их на сокет
func (gs *GlobalStorage) GetUsers(ctx *gin.Context) {
	body := &structs.User{}
	if err := ctx.BindJSON(&body); err != nil || body == nil {
		fmt.Printf("Ошибка в получении данных от сервера в GetUsers: %v\n", err)
		return
	}
	body.Messenger = ctx.Param("msgr")
	conn, err := socket.GetWSConnectClient("ws://localhost:8080/socket/cs") //Проверка сокета на запись
	if err != nil {
		log.Printf("Ошибка в подключении %v\n", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Ошибка в закрытии соединения: %v", err)
			return
		}
	}()
	go socket.WriteSocket(conn, gs.UserAdd(body))
	ctx.JSON(http.StatusAccepted, &body)
}

func (gs *GlobalStorage) GetData(ctx *gin.Context) {
	log.Println("GetData работает")
	conn, err := socket.GetWSConnectClient("ws://localhost:8080/socket/cs")
	if err != nil {
		log.Printf("Ошибка в подключении %v\n", err)
		return
	}
	body := &structs.User{}
	if err := ctx.BindJSON(&body); err != nil {
		fmt.Printf("Ошибка в получении данных от сервера в GetData: %v\n", err)
		return
	}
	go socket.WriteSocket(conn, gs.UserAdd(body))
	ctx.JSON(http.StatusAccepted, &body)
}

// PostData взятие данных с ReadSocketServer и отправка их на Клиента
func (gs *GlobalStorage) PostData(message []byte, url string) int {
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	config.SetHeaderPost(req)
	if err != nil {
		fmt.Printf("Ошибка в запросе: %v\n", err)
		return http.StatusBadRequest
	}
test:
	for {
		resp, err := client.Do(req)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Printf("Ошибка в подключении: %v\n", err)
			goto test
		}
		log.Printf("Отправлено от PostData по запросу: %v\n", string(message))
		if err := resp.Body.Close(); err != nil {
			log.Printf("Ошибка в закрытии ответа: %v\n", err)
			return http.StatusRequestTimeout
		}
	}
}

func (gs *GlobalStorage) Server(c *gin.Context) {
	fmt.Println("Сервер работает")
	var m *structs.User
	client := http.Client{}
	if err := c.BindJSON(&m); err != nil {
		log.Printf("Ошибка в получении JSON: %v", err)
		return
	}
	data, err := json.Marshal(&m)
	if err != nil {
		log.Printf("Ошибка в преобразовании JSON: %v", err)
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/users/getData", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Ошибка в закрытии ответа: %v", err)
			return
		}
	}()
}

func (gs *GlobalStorage) Client(c *gin.Context) {
	time.Sleep(1 * time.Second)
	fmt.Println("\nСервер работает")
	var m *structs.User
	client := http.Client{}
	if err := c.BindJSON(&m); err != nil {
		log.Printf("Ошибка в получении JSON: %v", err)
		return
	}
	data, err := json.Marshal(&m)
	if err != nil {
		log.Printf("Ошибка в преобразовании JSON: %v", err)
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/users/getUser/a", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Ошибка в закрытии ответа: %v", err)
			return
		}
	}()
}
