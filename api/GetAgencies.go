package api

import (
	"strconv"

	"github.com/OpenTransports/Paris/agencies"
	"github.com/OpenTransports/Paris/models"
	"github.com/go-siris/siris/context"
)

// GetAgencies - /api/agencies?latitude=...&longitude=...
// Send the agencies aroud the passed position or all agencies if no position is passed
// @formParam latitude : optional, the latitude around where to search, default is 0
// @formParam longitude : optional, the longitude around where to search, default is 0
func GetAgencies(ctx context.Context) {
	// Recover from potential panic in agencies
	defer func() {
		if r := recover(); r != nil {
			ctx.Application().Log(
				"Panic asking if agencies contains position (%v)\n	==> %v",
				ctx.FormValue("latitude"),
				ctx.FormValue("longitude"),
				r,
			)
			_, err := ctx.JSON(make([]models.ITransport, 0))
			if err != nil {
				ctx.Application().Log("Error writting answer in /api/agencies\n	==> %v", err)
			}
		}
	}()
	// Get position in params
	// Parse them to floats
	// Ignore errors because it default to 0
	latitude, _ := strconv.ParseFloat(ctx.FormValue("latitude"), 64)
	longitude, _ := strconv.ParseFloat(ctx.FormValue("longitude"), 64)
	// Return agencies that covers the position
	_, err := ctx.JSON(agencies.Containing(&models.Position{
		Latitude:  latitude,
		Longitude: longitude,
	}))
	// Log the error if any
	if err != nil {
		ctx.Application().Log("Error writting answer in /api/agencies\n	==> %v", err)
	}
}
