package integration_test

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/avito-pvz/internal/entity"
	. "github.com/Eun/go-hit"
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

	err := Do(
		Post(basePath+"/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(body),
		Expect().Status().Equal(http.StatusCreated),
		Store().Response().Body().JSON().JQ(".access_token").In(&accessToken),
	)

	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	return accessToken, nil
}
