package integration_test

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/avito-pvz/internal/entity"
	hit "github.com/Eun/go-hit"
)

const (
	moderatorEmail = "moderator@example.com"
	employeeEmail  = "employee@example.com"
	password       = "testpass"

	defaultAttempts = 20
	host            = "app:8080"
	healthPath      = "http://" + host + "/health"
	basePath        = "http://app:8080"
)

func Login(role string) (string, error) {
	var accessToken string

	email := moderatorEmail
	if role == string(entity.RoleEmployee) {
		email = employeeEmail
	}

	body := map[string]string{
		"email":    email,
		"password": password,
		"role":     role,
	}

	err := hit.Do(
		hit.Post(basePath+"/register"),
		hit.Send().Headers("Content-Type").Add("application/json"),
		hit.Send().Body().JSON(body),
		hit.Expect().Status().Equal(http.StatusCreated),
		hit.Store().Response().Body().JSON().JQ(".access_token").In(&accessToken),
	)

	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	return accessToken, nil
}
