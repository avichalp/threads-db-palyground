package main

import (
	"context"
	"fmt"
	"os"

	"tdb-client/todos_db"

	database "github.com/textileio/go-threads/db"
)

func main() {
	taskDescription := os.Args[1]
	keyPath := os.Args[2]

	// If private key in not found in the path
	// new one is generated and saved.
	privateKey, err := todos_db.GetPrivateKey(keyPath)
	if err != nil {
		panic(err)
	}

	// Get DB client using user's PrivKey
	db, threadID, err := todos_db.GetDBClient(privateKey)
	if err != nil {
		panic(err)
	}

	// TODO: Only let Admin create Schema
	todos_db.CreateSchema(db, threadID)
	if err != nil {
		panic(err)
	}

	// Get current user's token
	userToken, err := todos_db.GetUserToken(db, privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("User's Signed Thread Token", userToken)

	// TODO: Add user's token in its struct
	_, err = todos_db.AddNewToDoItem(db, threadID, taskDescription, "alice", string(userToken))

	// _, err = todos_db.AddNewToDoItem(db, threadID, taskDescription, "bob")
	if err != nil {
		panic(err)
	}

	// Find Bob's task with a query
	/* query := database.Where("person.name").Eq("bob")
	results, err := db.Find(context.Background(), threadID, "Todos", query, &todos_db.TodoItem{})
	if err != nil {
		panic(err)
	}
	item := results.([]*todos_db.TodoItem)[0]
	fmt.Println("Bob's task:", item.Description)
	*/
	// Alice's task
	query := database.Where("person.name").Eq("alice")
	results, err := db.Find(context.Background(), threadID, "Todos", query, &todos_db.TodoItem{})
	if err != nil {
		panic(err)
	}

	item := results.([]*todos_db.TodoItem)[0]
	fmt.Println("Alice's task:", item.Description)
	fmt.Println("Alice's signed token:", item.CreatedBy.Token)

}
