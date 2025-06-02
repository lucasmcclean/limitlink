package link

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongo(ctx context.Context) (*MongoDB, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, errors.New("MONGO_URI is required\n")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		return nil, errors.New("MONGO_DB_NAME is required\n")
	}

	collName := os.Getenv("MONGO_COLLECTION")
	if collName == "" {
		return nil, errors.New("MONGO_COLLECTION is required\n")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collName)

	if err := createIndexes(ctx, collection); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDB{
		client:     client,
		collection: collection,
	}, nil
}

func (db *MongoDB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

func (db *MongoDB) Create(ctx context.Context, link *Link) error {
	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now
	_, err := db.collection.InsertOne(ctx, link)
	return err
}

func (db *MongoDB) GetBySlug(ctx context.Context, slug string) (*Link, error) {
	var result Link
	err := db.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&result)
	if err != nil {
		return nil, err
	}
	err = Validate(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (db *MongoDB) GetAndInc(ctx context.Context, slug string) (*Link, error) {
	var updated Link
	err := db.collection.FindOneAndUpdate(
		ctx,
		bson.M{"slug": slug},
		bson.M{"$inc": bson.M{"hit_count": 1}},
	).Decode(&updated)
	if err != nil {
		return nil, err
	}
	err = Validate(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, err
}

func (db *MongoDB) GetByToken(ctx context.Context, token string) (*Link, error) {
	var result Link
	err := db.collection.FindOne(ctx, bson.M{"admin_token": token}).Decode(&result)
	return &result, err
}

func (db *MongoDB) DeleteByToken(ctx context.Context, token string) error {
	_, err := db.collection.DeleteOne(ctx, bson.M{"admin_token": token})
	return err
}

func (db *MongoDB) UpdateByToken(ctx context.Context, token string, updated *Link) error {
	updated.UpdatedAt = time.Now()
	_, err := db.collection.UpdateOne(
		ctx,
		bson.M{"admin_token": token},
		bson.M{"$set": updated},
	)
	return err
}

func createIndexes(ctx context.Context, coll *mongo.Collection) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("slug_idx"),
		},
		{
			Keys:    bson.D{{Key: "admin_token", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("admin_token_idx"),
		},
	}

	names, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	log.Printf("Created indexes: %v\n", names)
	return nil
}
