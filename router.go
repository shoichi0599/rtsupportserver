package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"net/http"
)

type Handler func(*Client, interface{})

// Switch protocols (http -> websocket)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow access from any origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Router struct {
	rules   map[string]Handler
	session *r.Session
}

func NewRouter(session *r.Session) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: session,
	}
}

//                                      handler func() string
//                                              -----  ------
//                                         Param Type  Return Type
//                                      => Create 'Handler' type as func for re-use
//                                         handler func() string -> handler Handler
// This method is to pass a message routing rule to our router.
func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	client := NewClient(socket, e.FindHandler, e.session)
	// 'defer client.Close()' is called before exiting ServeHTTP() function
	defer client.Close()
	// Since each of these methods should be running independently they'll need to be in separate go routines
	go client.Write()
	client.Read()
}
