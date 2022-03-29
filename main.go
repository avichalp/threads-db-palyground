package main

import (
	"fmt"
	"os"

	"tdb-client/todos_db"

	"github.com/spf13/cobra"
)

func start(name, keyPath, taskDescription string) {

	fmt.Println(name, keyPath, taskDescription)

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
	creator := todos_db.Person{
		ID:    "",
		Name:  name,
		Token: string(userToken),
	}
	_, err = todos_db.AddNewToDoItem(db, threadID, taskDescription, creator)

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

	/*
		query := database.Where("person.name").Eq("alice")
		results, err := db.Find(context.Background(), threadID, "Todos", query, &todos_db.TodoItem{})
		if err != nil {
			panic(err)
		}

		item := results.([]*todos_db.TodoItem)[0]
		fmt.Println("Alice's task:", item.Description)
		fmt.Println("Alice's signed token:", item.CreatedBy.Token) */
}

func main() {
	var name string
	var privKeyPath string
	var task string

	rootCmd := &cobra.Command{
		Use: "app",
	}
	startCmd := &cobra.Command{
		Use:   "todo",
		Short: "Add task",
		Long:  "Add task",
		Run: func(cmd *cobra.Command, args []string) {
			start(name, privKeyPath, task)
		},
	}

	startCmd.Flags().StringVarP(&name, "name", "n", "", "creator name")
	startCmd.Flags().StringVarP(&privKeyPath, "key", "f", "", "path to private key")
	startCmd.Flags().StringVarP(&task, "task", "t", "", "descrition of the task")

	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
