package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Client struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Enabled             bool   `json:"enabled"`
	Address             string `json:"address"`
	PublicKey           string `json:"public_key"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	DownloadableConfig  bool   `json:"downloadable_config"`
	PersistentKeepalive string `json:"persistent_keepalive"`
	LatestHandshakeAt   string `json:"latest_handshake_at"`
	TransferRX          int    `json:"transfer_rx"`
	TransferTX          int    `json:"transfer_tx"`
}

func init() {
	f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {

	key, err := login()
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := addClient("TestName", key)
	if err != nil {
		return
	}

	log.Println("Куки", key)

	log.Println("Клиент", client)

	list, err := clientList(key)
	if err != nil {
		return
	}
	fmt.Println(list)

	download, err := fileDownload(key)
	if err != nil {
		return
	}

	fmt.Println(download)
	log.Println("response download ", download)

	//start
	//Оплатить подписку
	// Ссылка на оплату -> Оплатил

	// Пришел платеж
	//login  +
	//addclient  +
	// get list client +
	// Client.Id   +
	// GET http://89.19.212.188:51821/api/wireguard/client/b9917b5a-1efc-4398-ab5b-3308ff24b7fa/configuration  +

	// Отправить файл с конфигурацией в телегам чат с пользователем.
}

func login() (*http.Cookie, error) {
	password := struct {
		Password string `json:"password"`
	}{Password: os.Getenv("CORRECT_PASS")}
	bytes, _ := json.Marshal(password)
	response, _ := http.Post(os.Getenv("WIREGUARD_LOGIN"), "application/json", strings.NewReader(string(bytes)))

	if response.StatusCode != 200 {
		return nil, errors.New("login failed")
	}

	auth := response.Cookies()[0]
	return auth, nil
}

func addClient(name string, key *http.Cookie) (bool, error) {
	addClientBody := struct {
		Name string `json:"name"`
	}{Name: name}

	addClientBodyToBytes, _ := json.Marshal(addClientBody)

	response, err := doRequestWireGuard("POST", os.Getenv("WIREGUARD_CLIENT"), addClientBodyToBytes, key)
	if err != nil {
		return false, errors.New("add client failed")
	}
	if response.StatusCode != 200 {
		return false, errors.New("add client failed")
	}
	return true, nil
}

func doRequestWireGuard(method string, url string, body []byte, key *http.Cookie) (*http.Response, error) {
	if method == "POST" {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(key)

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		body, _ = io.ReadAll(response.Body)

		return response, nil
	}
	if method == "GET" {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(key)
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
	return nil, errors.New("method not allowed")
}

func clientList(key *http.Cookie) (string, error) {

	response, err := doRequestWireGuard("GET", os.Getenv("WIREGUARD_CLIENT"), nil, key)
	if err != nil {

	}

	body, _ := io.ReadAll(response.Body)
	var result []Client
	err = json.Unmarshal(body, &result)
	lastClient := result[len(result)-1]
	lastClientID := lastClient.ID
	if err != nil {
		log.Fatal(err)
	}
	return lastClientID, nil
}

func fileDownload(key *http.Cookie) (*http.Response, error) {
	id, err := clientList(key)
	if err != nil {
		return nil, err
	}
	response, err := doRequestWireGuard("GET", os.Getenv("WIREGUARD_CLIENT")+id+"/configuration", nil, key)
	if err != nil {
	}
	return response, err
}
