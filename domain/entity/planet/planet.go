package planet

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Climate     string             `bson:"climate"`
	Terrain     string             `bson:"terrain"`
	Apparitions int32              `bson:"apparitions"`
}

func (p Planet) Validate() error {
	if p.Name == "" {
		return errors.New("invalid name param")
	}

	return nil
}
