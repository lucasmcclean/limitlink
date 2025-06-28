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

func (db *MongoDB) PatchByToken(ctx context.Context, token string, patch *link.PatchLink) error {
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

	_, err := db.collection.UpdateOne(ctx, bson.M{"admin_token": token}, updateDoc)
	return err
}
