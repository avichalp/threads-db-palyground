package todos_db

import (
	"context"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alecthomas/jsonschema"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/textileio/go-threads/api/client"
	"github.com/textileio/go-threads/core/thread"
	"google.golang.org/grpc"

	database "github.com/textileio/go-threads/db"
)

type Person struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type TodoItem struct {
	ID          string `json:"_id"`
	Description string `json:"description"`
	CreatedBy   Person `json:"person"`
	CreatedAt   int    `json:"created_at"`
}

func GetUserToken(db *client.Client, privateKey crypto.PrivKey) (thread.Token, error) {
	// Create a new ID by signing with the privateKey
	myIdentity := thread.NewLibp2pIdentity(privateKey)
	fmt.Println("MY PUB as base32 str", fmt.Sprintln(myIdentity.GetPublic()))

	threadToken, err := db.GetToken(context.Background(), myIdentity)
	if err != nil {
		return "", err
	}

	return threadToken, nil
}

func GetPrivateKey(keyPath string) (crypto.PrivKey, error) {
	var privateKey crypto.PrivKey
	var err error

	if keyPath == "" {
		privateKey, _, err = crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}

		privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		}
		err = ioutil.WriteFile("privkey", pem.EncodeToMemory(block), 0644)
		if err != nil {
			return nil, err
		}
	} else {
		b, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}

		firstBlock, _ := pem.Decode(b)
		privateKey, err = crypto.UnmarshalPrivateKey(firstBlock.Bytes)
		if err != nil {
			return nil, err
		}
	}

	return privateKey, nil
}

func GetDBClient(privateKey crypto.PrivKey) (*client.Client, thread.ID, error) {
	db, err := client.NewClient("127.0.0.1:6006", grpc.WithInsecure())
	if err != nil {
		return nil, "", err
	}

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

func AddNewToDoItem(db *client.Client, threadID thread.ID, task string, creator Person) ([]string, error) {
	todoItem := TodoItem{
		Description: task,
		CreatedBy:   creator,
		CreatedAt:   int(time.Now().UnixNano()),
	}

	ids, err := db.Create(context.Background(), threadID, "Todos", client.Instances{todoItem})
	if err != nil {
		return nil, err
	}

	return ids, nil
}
