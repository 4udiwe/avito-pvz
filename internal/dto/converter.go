package dto

import (
	"github.com/4udiwe/avito-pvz/internal/entity"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

//go:generate go tool oapi-codegen --config=dto.gen.yaml ../../api/swagger.yaml

func EntityPointToDTO(e *entity.Point) *PVZ {
	id := openapi_types.UUID(e.ID)
	return &PVZ{
		Id:               &id,
		City:             PVZCity(e.City),
		RegistrationDate: &e.CreatedAt,
	}
}

func EntityReceptionToDTO(e *entity.Reception) *Reception {
	id := openapi_types.UUID(e.ID)
	pointID := openapi_types.UUID(e.PointID)
	return &Reception{
		Id:       &id,
		PvzId:    pointID,
		DateTime: e.CreatedAt,
		Status:   ReceptionStatus(e.Status),
	}
}

func EntityProductToDTO(e *entity.Product) *Product {
	id := openapi_types.UUID(e.ID)
	receptionID := openapi_types.UUID(e.ReceptionID)
	return &Product{
		Id:          &id,
		ReceptionId: receptionID,
		DateTime:    &e.CreatedAt,
		Type:        ProductType(e.Type),
	}
}
