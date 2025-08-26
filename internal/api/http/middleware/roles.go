package middleware

import (
	"net/http"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles ...entity.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetUserFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			hasAccess := false
			for _, allowedRole := range allowedRoles {
				if claims.Role == allowedRole {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied")
			}

			return next(c)
		}
	}
}

func ModderatorOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(entity.RoleModerator)(next)
}

func EmployeeOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(entity.RoleEmployee)(next)
}

func EmployeeAndModerator(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(entity.RoleEmployee, entity.RoleModerator)(next)
}
