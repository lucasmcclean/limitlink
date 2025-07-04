package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lucasmcclean/limitlink/link"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Links wraps the "links" collection and implements the link.Repository
// interface.
type Links struct {
	collection *mongo.Collection
}

// Links returns a new Links wrapper for the store's "links" collection.
func (store *Store) Links() *Links {
	return &Links{
		store.db.Collection("links"),
	}
}

// ScheduleRemoval sets up a TTL index on the "admin_expires_at" field to
// automatically delete documents after the specified expiration time.
func (l *Links) EnsureTTLIndex(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys: bson.D{{Key: "admin_expires_at", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(0).
			SetName("adminExpiresAtTTL"),
	}

	_, err := l.collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return fmt.Errorf("failed to create TTL index: %w", err)
	}

	log.Println("TTL index on 'admin_expires_at' ensured.")
	return nil
}

// Create inserts a new link document into the collection.
func (l *Links) Create(ctx context.Context, vLink *link.Validated) error {
	_, err := l.collection.InsertOne(ctx, vLink.Link())
	return err
}

// GetBySlug retrieves a link document by its slug.
func (l *Links) GetBySlug(ctx context.Context, slug string) (*link.Link, error) {
	var result link.Link
	err := l.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	// TODO: Make sure link is publicly available
	return &result, err
}

// IncBySlug atomically increments the hit counter for the link with the given slug.
func (l *Links) IncBySlug(ctx context.Context, slug string) error {
	_, err := l.collection.UpdateOne(
		ctx,
		bson.M{"slug": slug},
		bson.M{"$inc": bson.M{"hit_count": 1}},
	)
	return err
}

// GetByToken retrieves a link document by its admin token.
func (l *Links) GetByToken(ctx context.Context, token string) (*link.Link, error) {
	var result link.Link
	err := l.collection.FindOne(ctx, bson.M{"admin_token": token}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &result, err
}

// DeleteByToken deletes a link document by its admin token.
func (l *Links) DeleteByToken(ctx context.Context, token string) error {
	_, err := l.collection.DeleteOne(ctx, bson.M{"admin_token": token})
	return err
}

// PatchByToken updates a link document by its admin token using a PatchLink struct.
func (l *Links) PatchByToken(ctx context.Context, token string, vPatch *link.ValidatedPatch) error {
	patch := vPatch.Patch()

	setFields := bson.M{
		"updated_at": patch.UpdatedAt,
	}
	unsetFields := bson.M{}

	if patch.ExpiresAt != nil {
		setFields["expires_at"] = *patch.ExpiresAt
	}

	if patch.MaxHits.Remove {
		unsetFields["max_hits"] = ""
	} else if patch.MaxHits.Value != nil {
		setFields["max_hits"] = *patch.MaxHits.Value
	}

	if patch.ValidFrom.Remove {
		unsetFields["valid_from"] = ""
	} else if patch.ValidFrom.Value != nil {
		setFields["valid_from"] = *patch.ValidFrom.Value
	}

	if patch.PasswordHash.Remove {
		unsetFields["password_hash"] = ""
	} else if patch.PasswordHash.Value != nil {
		setFields["password_hash"] = *patch.PasswordHash.Value
	}

	updateDoc := bson.M{}
	if len(setFields) > 0 {
		updateDoc["$set"] = setFields
	}
	if len(unsetFields) > 0 {
		updateDoc["$unset"] = unsetFields
	}

	if len(updateDoc) == 0 {
		return nil
	}

	_, err := l.collection.UpdateOne(ctx, bson.M{"admin_token": token}, updateDoc)
	return err
}
