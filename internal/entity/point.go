package entity

import (
	"time"

	"github.com/google/uuid"
)

type Point struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	City      string    `db:"city"`
}

type PointFullInfo struct {
	Point      Point
	Receptions []ReceptionWithProducts
}

type ReceptionWithProducts struct {
	Reception Reception
	Products  []Product
}
