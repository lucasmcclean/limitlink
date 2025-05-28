package link

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongo() (*MongoDB, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is required")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGO_DB_NAME is required")
	}

	collName := os.Getenv("MONGO_COLLECTION")
	if collName == "" {
		log.Fatal("MONGO_COLLECTION is required")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collName)
	return &MongoDB{collection: collection}, nil
}

func (db *MongoDB) Create(ctx context.Context, link *Link) error {
	link.CreatedAt = time.Now()
	link.UpdatedAt = time.Now()
	_, err := db.collection.InsertOne(ctx, link)
	return err
}

func (db *MongoDB) GetBySlug(ctx context.Context, slug string) (*Link, error) {
	var result Link
	err := db.collection.FindOne(ctx, map[string]interface{}{"slug": slug}).Decode(&result)
	return &result, err
}

func (db *MongoDB) GetAndInc(ctx context.Context, slug string) (*Link, error) {
	var updated Link
	err := db.collection.FindOneAndUpdate(
		ctx,
		map[string]interface{}{"slug": slug},
		map[string]interface{}{"$inc": map[string]int{"hit_count": 1}},
	).Decode(&updated)
	return &updated, err
}

func (db *MongoDB) GetByToken(ctx context.Context, token string) (*Link, error) {
	var result Link
	err := db.collection.FindOne(ctx, map[string]interface{}{"admin_token": token}).Decode(&result)
	return &result, err
}

func (db *MongoDB) DeleteByToken(ctx context.Context, token string) error {
	_, err := db.collection.DeleteOne(ctx, map[string]interface{}{"admin_token": token})
	return err
}

func (db *MongoDB) UpdateByToken(ctx context.Context, token string, updated *Link) error {
	updated.UpdatedAt = time.Now()
	_, err := db.collection.UpdateOne(
		ctx,
		map[string]interface{}{"admin_token": token},
		map[string]interface{}{"$set": updated},
	)
	return err
}
