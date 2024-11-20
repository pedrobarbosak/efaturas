package entity

import (
	"time"
)

type Entity struct {
	ID        int64 `bson:"_id" validate:"required"`
	CreatedAt int64 `validate:"required"`
	UpdatedAt int64
	DeletedAt int64
}

func New(id int64) *Entity {
	return &Entity{
		ID:        id,
		CreatedAt: time.Now().Unix(),
	}
}
