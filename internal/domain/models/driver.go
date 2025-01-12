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
	IsOnRide  bool               `bson:"is_on_ride"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Location struct {
	Coordinates []float64 `bson:"coordinates"`
}
