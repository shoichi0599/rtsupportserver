# Overview
Real Time Support (rtsupport) application on server side
This project is what I learned from the lecture (https://www.udemy.com/course/realtime-apps-with-reactjs-golang-rethinkdb/) 

Ref:
- https://github.com/knowthen/rtsupportclient
- https://github.com/knowthen/rtsupportserver

More Feature Suggestions from the lecture:
- Clientside Error Messages
- Editing / Deleting Messages
- Private Messaging
- Authentication with JWT's


# Dependencies
### RethinkDB-go - RethinkDB Driver for Go
github: https://github.com/rethinkdb/rethinkdb-go
```
# Installation
go get gopkg.in/rethinkdb/rethinkdb-go.v5
```

### Gorilla WebSocket
github: https://github.com/gorilla/websocket
```
# Installation
go get github.com/gorilla/websocket
```

### mapstructure
github: https://github.com/mitchellh/mapstructure
```
# Installation
go get github.com/mitchellh/mapstructure
```

### Goroutines
```go
// lightweight thread
go func() {...}
```
#### Communicate across Goroutine
Avoid Goroutines modifying shared data at same time.
Writing websocket in 2 Goroutines
```go
// These could lead to race condition.
// in subscribeChannel() goroutine
socket.WriteJson(message)
// in handler() goroutine
socket.WriteJson(message)
```
Popular saying
```
Do not communicate by sharing memory,
Shanre memory by communicating.
```

### Golang Channel
How you safely share data between Goroutines
```
Goroutine 1 | Goroutine 2
            |
       -----------
 data | <-(data)- | data (share)
       -----------
            |
            |
```
Example:
Channels provide a say way to pass values between Goroutines
```
Goroutine(main)       |       Goroutine(func)
                      |
                 -----------
 msg := "Hello" | <-(data)- | ("Hello")
                 -----------
                   msgChan
                      |
                      |
```
```go
package main

import "fmt"

func main()  {
	msgChan := make(chan string)
	go func() {
		msgChan <- "Hello"
	}()
	msg := <- msgChan
	fmt.Println(msg)
}
```

### RethinkDB
Install: https://rethinkdb.com/docs/install/<br>
Using Homebrew
```
brew update && brew install rethinkdb
```

### REQL - RethinkDB Query Language
REQL is embedded in programming languages
Create database and tables on http://localhost:8080/#dataexplorer
```sql
-- Create database
r.dbCreate('rtsupport')
#{
#  "config_changes": [
#    {
#      "new_val": {
#      "id":  "4198f4f5-e783-4e37-a400-47b85734033a" ,
#      "name":  "rtsupport"
#    } ,
#    "old_val": null
#    }
#  ] ,
#  "dbs_created": 1
#}

-- Create table
r.db('rtsupport').tableCreate('channel')
#{
#  "config_changes": [
#    {
#      "new_val": {
#        "db":  "rtsupport" ,
#        "durability":  "hard" ,
#        "id":  "625cb855-3fb2-4ceb-9550-16c457d636df" ,
#        "indexes": [ ],
#        "name":  "channel" ,
#        "primary_key":  "id" ,
#        "shards": [
#          {
#            "nonvoting_replicas": [ ],
#            "primary_replica":  "Shoichis_MacBook_Pro_local_e86" ,
#            "replicas": [
#              "MacBook_Pro_local_e86"
#            ]
#          }
#        ] ,
#        "write_acks":  "majority"
#      } ,
#      "old_val": null
#    }
#  ] ,
#  "tables_created": 1
#}

r.db('rtsupport').tableCreate('user')
...

r.db('rtsupport').tableCreate('message')
...

r.db('rtsupport').tableList()
# [
#   "channel" ,
#   "message" ,
#   "user"
# ]

-- 'channel' table
-- Insert
r.db('rtsupport').table('channel')
  .insert(
    { name: 'Hardware Support' }
  )

#{
#  "deleted": 0 ,
#  "errors": 0 ,
#  "generated_keys": [
#  "0bf48bc9-6237-4d45-a0ec-47c1cdd16f86"
#  ] ,
#  "inserted": 1 ,
#  "replaced": 0 ,
#  "skipped": 0 ,
#  "unchanged": 0
#}

-- Select all
r.db('rtsupport').table('channel')

-- Create index
r.db('rtsupport').table('channel')
  .indexCreate('name')
#{
# "created": 1
#}

-- 'user' table
-- Insert
r.db('rtsupport').table('user')
  .insert(
    { name: 'anonymous' }
  )
  
-- Select by
r.db('rtsupport').table('user').get('fd6f683a-c762-41fe-9580-aeabb6a69b86')

-- Update
r.db('rtsupport').table('user').get('fd6f683a-c762-41fe-9580-aeabb6a69b86')
  .update({ name: 'Test' })

-- Create index
r.db('rtsupport').table('user')
  .indexCreate('name')

-- Delete
r.db('rtsupport').table('user')
  .get('fd6f683a-c762-41fe-9580-aeabb6a69b86')
  .delete();

-- 'message' table
r.db('rtsupport').table('message')
  .insert({
    author: 'Test',
    createdAt: r.now(),
    body: 'I need some help...',
    channelId: '0bf48bc9-6237-4d45-a0ec-47c1cdd16f86'
  })
  
r.db('rtsupport').table('message')
  .indexCreate('createdAt')

r.db('rtsupport').table('message')
  .indexCreate('channelId')
```

### REQL Changefeed
- Optionally Provides Initial Query Records
- Streams any new changes in Realtime

Changefeed Result
```
{
  "new_val": {...},
  "old_val": {...}
}
```
```sql
r.db('rtsupport').table('channel')
  .changes({includeInitial: true})
#{
#  "new_val": {
#  "id":  "0bf48bc9-6237-4d45-a0ec-47c1cdd16f86" ,
#  "name":  "Hardware Support"
#  }
#}

-- Insert data
-- r.db('rtsupport').table('channel')
--  .insert({name: 'Software Support'})
#{
#  "new_val": {
#  "id":  "02bd6a90-9735-4cc6-9039-e2129c0fb120" ,
#  "name":  "Software Support"
#  } ,
#  "old_val": null
#  } {
#  "new_val": {
#  "id":  "0bf48bc9-6237-4d45-a0ec-47c1cdd16f86" ,
#  "name":  "Hardware Support"
#  }
#}

-- Update data
-- r.db('rtsupport').table('channel').get('02bd6a90-9735-4cc6-9039-e2129c0fb120')
--   .update({name: 'Critical Software Support'})
#{
#  "new_val": {
#  "id":  "02bd6a90-9735-4cc6-9039-e2129c0fb120" ,
#  "name":  "Critical Software Support"
#  } ,
#  "old_val": {
#  "id":  "02bd6a90-9735-4cc6-9039-e2129c0fb120" ,
#  "name":  "Software Support"
#  }
#}
```
### REQL in GO
main.go
```go
package main

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type User struct {
	Id string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
		Database: "rtsupport",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	user := User {
		Name: "anonymous",
	}
	response, err := r.Table("user").
		Insert(user).
		RunWrite(session)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", response)
}
```
```
$ go run main.go 
rethinkdb.WriteResponse{Errors:0, Inserted:1, Updated:0, Unchanged:0, Replaced:0, Renamed:0, Skipped:0, Deleted:0, Created:0, DBsCreated:0, TablesCreated:0, Dropped:0, DBsDropped:0, TablesDropped:0, GeneratedKeys:[]string{"46f5193f-e1bc-4c2f-a717-f3ebe3de4485"}, FirstError:"", ConfigChanges:[]rethinkdb.ChangeResponse(nil), Changes:[]rethinkdb.ChangeResponse(nil)}

```

### Database Connection
```
   ________
  |        |
  | Main() | session := r.Connect(...)
  |________|
      ↓
   ________
  |        |
  | Router | router := NewRouter(session)
  |________|
      ↓
   ________
  |        |
  | Client | client := NewClient(session)
  |________|
```

- Subscription Handling
```
   _________      __________               ______________
  |         |    |          | Go Channel  |              |
  | Browser | -> |  Client  | --------->  | Subscription |
  |_________|    | (Struct) |             |   Handler    |
                 |__________|             |______________|

  Disconnect       "Stop"     -("Stop")->     "Stop"
                                          // Kill change feed
                                          cursor.Close()
                                          // Exit goroutines
                                          return
```
```go
package main

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"time"
)

func subscribe(session *r.Session, stop <-chan bool) {
	result := make(chan r.ChangeResponse)
	cursor, _ := r.Table("channel").
		Changes().
		Run(session)
	go func() {
		var change r.ChangeResponse
		for cursor.Next(&change) {
			// In actual app, send update to client
			//fmt.Printf("%#v\n", change.NewValue)
			result <- change
		}
		fmt.Println("exiting cursor goroutine")
	}()
	// Go select
	for {
		select {
		case change := <-result:
			fmt.Printf("%#v\n", change.NewValue)
		case <-stop:
			fmt.Println("closing cursor")
			cursor.Close()
			return
		}
	}
}

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "rtsupport",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	stop := make(chan bool)
	go subscribe(session, stop)
	// sleep to keep app running
	time.Sleep(time.Second * 5)
	fmt.Println("sending stop")
	stop <- true
	fmt.Println("browser closes... websocket closes")
	time.Sleep(time.Second * 10000)
}
```

### Go Select
```go
// Wait on multiple Channels
for {
	select {
	case val1 := <- chan1
	// do something with val1
	case val2 := <- chan2
	// do something with val2
	}
}
```

### Go Defere
```go
// If you call 'defer Cleanup()' at the top of the function
// the Cleanup() won't happen right away
// The Cleanup() function is queued up to get called before exiting the current function
func SomeFunction() {
	defer Cleanup()
	Step1()
	Step2()
	Step3()
	// Cleanup() called before exiting SomeFunction()
}
```