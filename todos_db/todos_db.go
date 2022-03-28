package todos_db

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/alecthomas/jsonschema"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/textileio/go-threads/api/client"
	"github.com/textileio/go-threads/core/thread"
	"google.golang.org/grpc"

	database "github.com/textileio/go-threads/db"
)

type Person struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type TodoItem struct {
	ID          string `json:"_id"`
	Description string `json:"description"`
	CreatedBy   Person `json:"person"`
	CreatedAt   int    `json:"created_at"`
}

func GetDBClient() (*client.Client, thread.ID, error) {
	db, err := client.NewClient("127.0.0.1:6006", grpc.WithInsecure())

	// Private key is kept locally
	privateKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, "", err
	}

	// Create a new ID by signing with the privateKey
	myIdentity := thread.NewLibp2pIdentity(privateKey)
	threadToken, err := db.GetToken(context.Background(), myIdentity)
	if err != nil {
		return nil, "", err
	}

	fmt.Println("Signed Thread Token", threadToken)

	threadID := thread.NewIDV1(thread.Raw, 32)
	err = db.NewDB(context.Background(), threadID)
	if err != nil {
		return nil, "", err
	}

	return db, threadID, nil
}

func CreateSchema(db *client.Client, threadID thread.ID) error {
	reflector := jsonschema.Reflector{}
	mySchema := reflector.Reflect(TodoItem{})

	err := db.NewCollection(context.Background(), threadID, database.CollectionConfig{
		Name:   "Todos",
		Schema: mySchema,
		Indexes: []database.Index{{
			Path:   "person.name", // Value matches json tags
			Unique: true,          // Create a unique index on "name"
		}},
	})

	if err != nil {
		return err
	}
	return nil
}

func AddNewToDoItem(db *client.Client, threadID thread.ID, task string, creator string) ([]string, error) {
	todoItem := TodoItem{
		Description: task,
		CreatedBy: Person{
			ID:   "",
			Name: creator,
		},
		CreatedAt: int(time.Now().UnixNano()),
	}

	ids, err := db.Create(context.Background(), threadID, "Todos", client.Instances{todoItem})
	if err != nil {
		return nil, err
	}

	return ids, nil
}
