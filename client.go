package main

import (
	"github.com/gorilla/websocket"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
)

type FindHandler func(string) (Handler, bool)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Responsibilities of the Client
// => Send & Receive Messages to Browser
type Client struct {
	send   chan Message
	socket *websocket.Conn
	// findHandler func(string) (Handler, bool)
	// Create a type like 'FindHandler' for 'func(string) (Handler, bool)' for readability
	findHandler  FindHandler
	session      *r.Session
	stopChannels map[int]chan bool
	id           string
	userName     string
}

func (c *Client) NewStopChannel(stopKey int) chan bool {
	// Guard Against Client Side Bug
	// Browser ==> 'channel subscribe' ==> subscribeChannel() ==> `stop := NewStopChannel(ChannelStop)`
	//                                                             -> adds stop channel to map in client
	//         ==> 'channel subscribe' ==> subscribeChannel() ==> `stop := NewStopChannel(ChannelStop)`
	//                                                             -> overwrites stop channel in client
	//                                                                -> goroutine leak
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		// we know we should call a function to handle the message
		// what function to call?
		// => create function to find a handler
		if handler, found := client.findHandler(message.Name); found {
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

// Capitalize 'w'
// to make it public or accessible outside its own package we'll capitalize the first letter W
func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.send)
}

func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	var user User
	user.Name = "anonymous"
	res, err := r.Table("user").Insert(user).RunWrite(session)
	if err != nil {
		log.Println(err.Error())
	}
	var id string
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	return &Client{
		send:         make(chan Message),
		socket:       socket,
		findHandler:  findHandler,
		session:      session,
		stopChannels: make(map[int]chan bool),
		id:           id,
		userName:     user.Name,
	}
}
