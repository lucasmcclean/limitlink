package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/lucasmcclean/limitlink/link"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrLinkUnavailable = errors.New("The requested link is unavailable.")

func (db *MongoDB) Create(ctx context.Context, link *link.Link) error {
	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now
	_, err := db.collection.InsertOne(ctx, link)
	return err
}

func (db *MongoDB) GetBySlug(ctx context.Context, slug string) (*link.Link, error) {
	var result link.Link
	err := db.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&result)
	if err != nil {
		return nil, err
	}
	available := result.IsAvailable()
	if !available {
		return nil, ErrLinkUnavailable
	}
	return &result, nil
}

func (db *MongoDB) IncBySlug(ctx context.Context, slug string) error {
	err := db.collection.FindOneAndUpdate(
		ctx,
		bson.M{"slug": slug},
		bson.M{"$inc": bson.M{"hit_count": 1}},
	).Err()
	return err
}

func (db *MongoDB) GetByToken(ctx context.Context, token string) (*link.Link, error) {
	var result link.Link
	err := db.collection.FindOne(ctx, bson.M{"admin_token": token}).Decode(&result)
	return &result, err
}

func (db *MongoDB) DeleteByToken(ctx context.Context, token string) error {
	_, err := db.collection.DeleteOne(ctx, bson.M{"admin_token": token})
	return err
}

func (db *MongoDB) UpdateByToken(ctx context.Context, token string, updated *link.Link) error {
	updated.UpdatedAt = time.Now()
	_, err := db.collection.UpdateOne(
		ctx,
		bson.M{"admin_token": token},
		bson.M{"$set": updated},
	)
	return err
}
