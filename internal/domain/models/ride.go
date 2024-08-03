package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	RideStatus string
)

const (
	RideStatusPending    RideStatus = "pending"
	RideStatusInProgress RideStatus = "in_progress"
	RideStatusCompleted  RideStatus = "completed"
)

var (
	RideStatusToString = map[RideStatus]string{
		RideStatusPending:    "pending",
		RideStatusInProgress: "in_progress",
		RideStatusCompleted:  "completed",
	}

	StringToRideStatus = map[string]RideStatus{
		"pending":     RideStatusPending,
		"in_progress": RideStatusInProgress,
		"completed":   RideStatusCompleted,
	}
)

type Ride struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Rider         User               `bson:"rider"`
	Driver        User               `bson:"driver"`
	StartLocation Location           `bson:"start_location"`
	DropLocation  Location           `bson:"drop_location"`
	Price         float64            `bson:"price"`
	Distance      float64            `bson:"distance"`
	Status        RideStatus         `bson:"status"`
	Verification  string             `bson:"verification"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type User struct {
	ID    string `bson:"id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}
