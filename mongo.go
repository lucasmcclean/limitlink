package main

import "go.mongodb.org/mongo-driver/mongo"

type MongoDB struct {
	collection *mongo.Collection
}
