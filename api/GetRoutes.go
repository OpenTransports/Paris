package api

import (
	"github.com/OpenTransports/Paris/agencies"
	"github.com/OpenTransports/lib-go/models"
	"github.com/go-siris/siris/context"
)

// GetRoutes - /api/agencies?latitude=...&longitude=...
// Send all the routes for all the agencies available
func GetRoutes(ctx context.Context) {
	// Get all the routes from all the agencies
	routes := make([]models.Route, 0)
	for _, agency := range agencies.AllAgencies {
		routes = append(routes, agency.Routes...)
	}

	_, err := ctx.JSON(routes)
	// Log the error if any
	if err != nil {
		ctx.Application().Logger().Errorf("Error writting answer in /api/agencies\n	==> %v", err)
	}
}
