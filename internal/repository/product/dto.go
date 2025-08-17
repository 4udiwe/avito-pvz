package repo_product

import (
	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

type CreateProduct struct {
	ReceptionID uuid.UUID
	Type        entity.ProductType
}
