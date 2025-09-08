//go:build integration

package integration_test

import (
	"net/http"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	. "github.com/Eun/go-hit"
	"github.com/google/uuid"
)

func TestReceptionFullCycle(t *testing.T) {
	const productsAmount = 50
	const city dto.PVZCity = "Москва"

	productTypes := []entity.ProductType{entity.ProductTypeClothes, entity.ProductTypeElectronics, entity.ProductTypeShoes}

	moderatorToken, err := Login(string(entity.RoleModerator))
	if err != nil {
		t.Fatal(err)
	}

	employeeToken, err := Login(string(entity.RoleEmployee))
	if err != nil {
		t.Fatal(err)
	}

	pointID, err := createPoint(moderatorToken, string(city))
	if err != nil {
		t.Fatal(err)
	}

	err = openReception(employeeToken, pointID)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < productsAmount; i++ {
		err = createProduct(employeeToken, pointID, productTypes[i%len(productTypes)])
		if err != nil {
			t.Fatal(err)
		}
	}

	err = closeReception(employeeToken, pointID)
	if err != nil {
		t.Fatal(err)
	}
}

func createPoint(moderatorToken string, city string) (uuid.UUID, error) {
	var id uuid.UUID

	body := map[string]string{"city": city}

	if err := Do(
		Post(basePath+"/pvz"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add("Bearer "+moderatorToken),
		Send().Body().JSON(body),
		Expect().Status().Equal(http.StatusCreated),
		Store().Response().Body().JSON().JQ(".id").In(&id),
	); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func openReception(employeeToken string, pointID uuid.UUID) error {
	body := map[string]string{"pvzId": pointID.String()}

	if err := Do(
		Post(basePath+"/receptions"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add("Bearer "+employeeToken),
		Send().Body().JSON(body),
		Expect().Status().Equal(http.StatusCreated),
	); err != nil {
		return err
	}
	return nil
}

func closeReception(employeeToken string, pointID uuid.UUID) error {
	if err := Do(
		Post(basePath+"/pvz/"+pointID.String()+"/close_last_reception"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add("Bearer "+employeeToken),
		Expect().Status().Equal(http.StatusAccepted),
	); err != nil {
		return err
	}
	return nil
}

func createProduct(employeeToken string, pointID uuid.UUID, productType entity.ProductType) error {
	body := map[string]string{"pvzId": pointID.String(), "type": string(productType)}

	if err := Do(
		Post(basePath+"/products"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add("Bearer "+employeeToken),
		Send().Body().JSON(body),
		Expect().Status().Equal(http.StatusCreated),
	); err != nil {
		return err
	}
	return nil
}
