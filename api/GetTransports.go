package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/OpenTransports/Paris/agencies"
	"github.com/OpenTransports/Paris/agencies/ratp"
	"github.com/OpenTransports/lib-go/models"
	"github.com/go-siris/siris/context"
)

// GetTransports - /api/transports?latitude=...&longitude=...&radius=...
// Send the transports aroud the passed position
// @formParam latitude : optional, the latitude around where to search, default is 0
// @formParam longitude : optional, the longitude around where to search, default is 0
// @formParam radius : optional, default is 200m
func GetTransports(ctx context.Context) {
	// Recover from potential panic in agencies
	defer func() {
		if r := recover(); r != nil {
			ctx.Application().Logger().Errorf("Panic asking for nearest transports to agencies\n	==> %v", r)
			_, err := ctx.JSON(make([]models.Transport, 0))
			if err != nil {
				ctx.Application().Logger().Errorf("Error writting answer in /api/transports\n	==> %v", err)
			}
		}
	}()
	// Get position in params
	// Parse them to floats
	// Ignore errors because it default to 0
	latitude, _ := strconv.ParseFloat(ctx.FormValue("latitude"), 64)
	longitude, _ := strconv.ParseFloat(ctx.FormValue("longitude"), 64)
	radius, _ := strconv.ParseInt(ctx.FormValue("radius"), 10, 64)
	// Create a Position object
	position := models.Position{
		Latitude:  latitude,
		Longitude: longitude,
	}
	// Set the radius to its default value if none is passed
	if radius == 0 {
		radius = 200
	}
	// Create transports array - max size is 300
	nearestTransports := []models.Transport{}
	// Get all transports near the passed position
	// Only ask agencies that cover the passed position
	for _, agency := range agencies.Containing(position) {
		transports := agency.TransportsNearPosition(position, int(radius))
		nearestTransports = append(nearestTransports, transports...)
	}
	// Get fresh infos for each transports
	errs := UpdateInfos(nearestTransports)
	for _, err := range errs {
		ctx.Application().Logger().Errorf("Error in /api/transports\n	==> %v", err)
	}
	// Return the result
	_, err := ctx.JSON(nearestTransports)
	// Log the error if any
	if err != nil {
		ctx.Application().Logger().Errorf("Error writting answer in /api/transports\n	==> %v", err)
	}
}

// UpdateInfos - Refresh information concerning transports
// @params transports : the transports list to update
// @return encontered errors
func UpdateInfos(transports []models.Transport) []error {
	if len(transports) == 0 {
		return nil
	}
	// Create a chanel
	errors := make(chan error, len(transports))
	gather := make(chan bool, len(transports))

	// For each transports, concurencialy update infos
	for i := range transports {

		go func(i int) {
			errC := make(chan error, 1)
			doneC := make(chan bool, 1)

			go func() {
				// Recover from potential panic in agencies
				defer func() {
					if r := recover(); r != nil {
						errC <- fmt.Errorf("Panic updating transport\n	==> %v\n	==> station: (%v)", r, transports[i])
					}
				}()
				err := ratp.UpdateTransportInfo(&transports[i])
				if err != nil {
					panic(err)
				} else {
					doneC <- true
				}
			}()

			select {
			case done := <-doneC:
				gather <- done
			case err := <-errC:
				errors <- err
			case <-time.After(7 * time.Second):
				errors <- fmt.Errorf("Time out after 3s\n	==> station: %v", transports[i])
			}

		}(i)

	}
	// Wait for each update to finish
	i := 0
	errs := make([]error, 0, len(transports))
	for range transports {
		select {
		case <-gather:
			continue
		case err := <-errors:
			errs = append(errs, err)
		}
		i++
		if i == len(transports) {
			break
		}
	}
	return errs
}
