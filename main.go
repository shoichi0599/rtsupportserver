package main

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"log"
	"net/http"
)

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "rtsupport",
	})
	if err != nil {
		log.Panic(err.Error())
	}
	router := NewRouter(session)

	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)
	router.Handle("channel unsubscribe", unsubscribeChannel)

	router.Handle("user edit", editUser)
	router.Handle("user subscribe", subscribeUser)
	router.Handle("user unsubscribe", unsubscribeUser)

	router.Handle("message add", addChannelMessage)
	router.Handle("message subscribe", subscribeChannelMessage)
	router.Handle("message unsubscribe", unsubscribeChannelMessage)

	//http.HandleFunc("/", handler)
	http.Handle("/", router)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		fmt.Println(err)
		return
	}
}
