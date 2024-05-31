package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Booking struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomId primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
}
