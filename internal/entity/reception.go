package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReceptionStatus string

const (
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
	ReceptionStatusClosed     ReceptionStatus = "close"
)

type Reception struct {
	ID        uuid.UUID       `db:"id"`
	PointID   uuid.UUID       `db:"point_id"`
	CreatedAt time.Time       `db:"created_at"`
	Status    ReceptionStatus `db:"status"`
}
