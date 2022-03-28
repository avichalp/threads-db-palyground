package main

import (
	"context"
	"fmt"

	"tdb-client/todos_db"

	database "github.com/textileio/go-threads/db"
)

func main() {

	db, threadID, err := todos_db.GetDBClient()

	todos_db.CreateSchema(db, threadID)
	if err != nil {
		panic(err)
	}

	_, err = todos_db.AddNewToDoItem(db, threadID, "buy milk", "alice")
	_, err = todos_db.AddNewToDoItem(db, threadID, "do laundry", "bob")
	if err != nil {
		panic(err)
	}

	// Find Bob's task with a query
	query := database.Where("person.name").Eq("bob")
	results, err := db.Find(context.Background(), threadID, "Todos", query, &todos_db.TodoItem{})
	if err != nil {
		panic(err)
	}
	item := results.([]*todos_db.TodoItem)[0]
	fmt.Println("Bob's task:", item.Description)

	// Alice's task
	query = database.Where("person.name").Eq("alice")
	results, err = db.Find(context.Background(), threadID, "Todos", query, &todos_db.TodoItem{})
	if err != nil {
		panic(err)
	}

	item = results.([]*todos_db.TodoItem)[0]
	fmt.Println("Alice's task:", item.Description)

}
