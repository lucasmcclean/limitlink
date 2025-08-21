package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store holds the MongoDB client and database reference.
type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

// New creates a new Store by connecting to MongoDB using environment variables.
// It returns an error if the environment variables are missing or if connection
// or ping fails.
func New(ctx context.Context) (*Store, error) {
	uri, uriPresent := os.LookupEnv("MONGO_URI")
	dbName, dbNamePresent := os.LookupEnv("MONGO_NAME")

	missing := make([]string, 0, 2)
	if !uriPresent {
		missing = append(missing, "MONGO_URI")
	}
	if !dbNamePresent {
		missing = append(missing, "MONGO_NAME")
	}
	if len(missing) != 0 {
		return nil, errors.New("missing one or more environment variable: " + strings.Join(missing, ", "))
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error verifying MongoDB connection: %w", err)
	}

	return &Store{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// Close disconnects the MongoDB client.
func (store *Store) Close(ctx context.Context) error {
	return store.client.Disconnect(ctx)
}
