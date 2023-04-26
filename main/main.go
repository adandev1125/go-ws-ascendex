package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func (s *APIClientStruct) Connection() error {
	var err error
	var response *http.Response

	if s.dialer == nil {
		return fmt.Errorf("connect with nil dialer")
	}

	url := url.URL{Scheme: "wss", Host: "ascendex.com", Path: "/0/api/pro/v1/stream"}

	s.conn, response, err = s.dialer.Dial(url.String(), nil)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

func (s *APIClientStruct) SubscribeToChannel(symbol string) error {
	if s.dialer == nil || s.conn == nil {
		log.Fatal("subscribe with nil dialer or connection")
	}

	var token string
	var asset string

	parts := strings.Split(symbol, "_")
	if len(parts) == 2 {
		token = parts[0]
		asset = parts[1]
	} else {
		return fmt.Errorf("invalid symbol format")
	}

	var subscribeMessage = SubscribeMessage{
		Op: "sub",
		Id: "go-apiclient-test",
		Ch: fmt.Sprintf("bbo:%s/%s", asset, token),
	}

	s.conn.WriteJSON(subscribeMessage)

	for {
		var message ReceiveMessage
		err := s.conn.ReadJSON(&message)

		if err != nil {
			return err
		}

		if message.M == "sub" && message.Id == subscribeMessage.Id && message.Ch == subscribeMessage.Ch && message.Code == 0 {
			log.Printf("subscribed to the channel: %s\n", message.Ch)
			break
		}
	}

	messages := make(chan BestOrderBook)

	go func() {
		s.ReadMessagesFromChannel(messages)
	}()

	go func() {
		s.WriteMessagesToChannel()
	}()

	for message := range messages {
		log.Printf("Received message: %#v", message)
	}

	return nil
}

func (s *APIClientStruct) ReadMessagesFromChannel(ch chan<- BestOrderBook) {
	for {
		var message BBOBlob
		err := s.conn.ReadJSON(&message)

		if err != nil {
			log.Fatal(err)
			continue
		}

		if message.M != "bbo" {
			continue
		}

		var bestOrderBook BestOrderBook
		err = message.Data.ToBestOrderBook(&bestOrderBook)
		if err != nil {
			log.Fatal(err)
			continue
		}
		ch <- bestOrderBook
	}
}

func (s *APIClientStruct) WriteMessagesToChannel() {
	for {
		var pingMessage = PingMessage{
			Op: "Ping",
		}
		s.conn.WriteJSON(pingMessage)

		time.Sleep(15 * time.Second)
	}
}

func (s *APIClientStruct) Disconnect() {
	s.conn.Close()
}

func ConnectionHandler(connectionMessage ReceiveMessage, apiClient APIClientStruct) {
	log.Printf("Connected")
	apiClient.SubscribeToChannel("USDT_BTC")
}

func main() {
	// create a Dialer

	var err error

	var apiClient = APIClientStruct{
		dialer: websocket.DefaultDialer,
	}

	err = apiClient.Connection()

	if err != nil {
		log.Fatal(err)
		return
	}

	defer apiClient.Disconnect()

	for {
		_, message, err := apiClient.conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}

		// log.Println(string(message))

		var receive ReceiveMessage

		err = json.Unmarshal(message, &receive)

		if err != nil {
			log.Fatal(err)
			continue
		}

		if receive.M == "connected" {
			ConnectionHandler(receive, apiClient)
			return
		}
	}
}
