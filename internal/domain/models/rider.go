package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rider struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	ContactNo string             `bson:"contact_no,omitempty"`
	Email     string             `bson:"email,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}
