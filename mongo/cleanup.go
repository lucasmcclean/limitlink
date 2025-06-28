package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const cleanupInterval = time.Second * 10

func (db *MongoDB) scheduleLinksCleanup(ctx context.Context) {
	ticker := time.NewTicker(cleanupInterval)

	go func() {
		for range ticker.C {
			if err := db.cleanupLinks(ctx); err != nil {
				log.Println("Error during cleanup:", err)
			}
		}
	}()
}

func (db *MongoDB) cleanupLinks(ctx context.Context) error {
	filter := bson.M{
		"$or": []bson.M{
			{"expires_at": bson.M{"$lt": time.Now()}},
			{
				"$and": []bson.M{
					{"max_hits": bson.M{"$ne": nil}},
					{"$expr": bson.M{"$gte": []any{"$hit_count", "$max_hits"}}},
				},
			},
		},
	}

	res, err := db.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	log.Printf("Deleted %d expired or max-hit links\n", res.DeletedCount)
	return nil
}
