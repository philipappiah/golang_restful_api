package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsFeed struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
}
