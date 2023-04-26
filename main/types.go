package main

import (
	"strconv"

	"github.com/gorilla/websocket"
)

type APIClient interface {
	/*
		Implement a websocket connection function
	*/
	Connection() error

	/*
		Implement a disconnect function from websocket
	*/
	Disconnect()

	/*
		Implement a function that will subscribe to updates
		of BBO for a given symbol

		The symbol must be of the form "TOKEN_ASSET"
		As an example "USDT_BTC" where USDT is TOKEN and BTC is ASSET

		You will need to convert the symbol in such a way that
		it complies with the exchange standard
	*/
	SubscribeToChannel(symbol string) error

	/*
		Implement a function that will write the data that
		we receive from the exchange websocket to the channel
	*/
	ReadMessagesFromChannel(ch chan<- BestOrderBook)

	/*
		Implement a function that will support connecting to a websocket
	*/
	WriteMessagesToChannel()
}

type APIClientStruct struct {
	dialer *websocket.Dialer
	conn   *websocket.Conn
}

// BestOrderBook struct
type BestOrderBook struct {
	Ask Order `json:"ask"` //asks.Price > any bids.Price
	Bid Order `json:"bid"`
}

// Order struct
type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

// Receiving Message
type ReceiveMessage struct {
	M    string `json:"m"`
	Type string `json:"type"`
	Hp   int    `json:"hp"`
	Id   string `json:"id"`
	Ch   string `json:"ch"`
	Code int    `json:"code"`
}

// Ping Message to keep connection alive
type PingMessage struct {
	Op string `json:"op"`
}

// Subscribe Message
type SubscribeMessage struct {
	Op string `json:"op"`
	Id string `json:"id"`
	Ch string `json:"ch"`
}

// BBO Blob struct
type BBOBlob struct {
	M      string      `json:"m"`
	Symbol string      `json:"symbol"`
	Data   BBOBlobData `json:"data"`
}

type BBOBlobData struct {
	Ts  int64     `json:"ts"`
	Bid [2]string `json:"bid"`
	Ask [2]string `json:"ask"`
}

func (b *BBOBlobData) ToBestOrderBook(bestOrderBook *BestOrderBook) error {

	askAmount, err := strconv.ParseFloat(b.Ask[0], 64)
	if err != nil {
		return err
	}

	askPrice, err := strconv.ParseFloat(b.Ask[1], 64)
	if err != nil {
		return err
	}

	bidAmount, err := strconv.ParseFloat(b.Bid[0], 64)
	if err != nil {
		return err
	}

	bidPrice, err := strconv.ParseFloat(b.Bid[1], 64)
	if err != nil {
		return err
	}

	bestOrderBook.Ask.Amount = askAmount
	bestOrderBook.Ask.Price = askPrice
	bestOrderBook.Bid.Amount = bidAmount
	bestOrderBook.Bid.Price = bidPrice

	return nil
}
