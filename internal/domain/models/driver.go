package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Driver struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	License   string             `bson:"license,omitempty"`
	ContactNo string             `bson:"contact_no,omitempty"`
	Email     string             `bson:"email,omitempty"`
	IsActive  bool               `bson:"is_active,omitempty"`
	Password  string             `bson:"password,omitempty"`
	Jwt       string             `bson:"jwt,omitempty"`
	Location  Location           `bson:"location,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

type Location struct {
	XCoordinate uint64 `bson:"x"`
	YCoordinate uint64 `bson:"y"`
}
