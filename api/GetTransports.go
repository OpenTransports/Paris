package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/OpenTransports/Paris/agencies"
	"github.com/OpenTransports/Paris/models"
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
			ctx.Application().Log("Panic asking for nearest transports to agencies\n	==> %v", r)
			_, err := ctx.JSON(make([]models.ITransport, 0))
			if err != nil {
				ctx.Application().Log("Error writting answer in /api/transports\n	==> %v", err)
			}
		}
	}()
	// Get position in params
	// Parse them to floats
	// Ignore errors because it default to 0
	latitude, _ := strconv.ParseFloat(ctx.FormValue("latitude"), 64)
	longitude, _ := strconv.ParseFloat(ctx.FormValue("longitude"), 64)
	radius, _ := strconv.ParseFloat(ctx.FormValue("radius"), 64)
	// Create a Position object
	position := &models.Position{
		Latitude:  latitude,
		Longitude: longitude,
	}
	// Set the radius to its default value if none is passed
	if radius == 0 {
		radius = 200.0
	}
	// Create transports array - max size is 300
	nearestTransports := []models.ITransport{}
	// Get all transports near the passed position
	// Only aske agencies that cover the passed position
	for _, a := range agencies.Containing(position) {
		transports, err := a.TransportsNearPosition(position, radius)
		if err != nil {
			ctx.Application().Log("Error in /api/transports\n	==> %v", err)
			continue
		}
		nearestTransports = append(nearestTransports, transports...)
	}
	// Get fresh infos for each transports
	errs := updateInfos(nearestTransports)
	for _, err := range errs {
		ctx.Application().Log("Error in /api/transports\n	==> %v", err)
	}
	// Return the result
	_, err := ctx.JSON(nearestTransports)
	// Log the error if any
	if err != nil {
		ctx.Application().Log("Error writting answer in /api/transports\n	==> %v", err)
	}
}

// Refresh information concerning transports
// @params transports : the transports list to update
func updateInfos(transports []models.ITransport) (err []error) {
	if len(transports) == 0 {
		return nil
	}
	// Create a chanel
	errors := make(chan error, len(transports))
	gather := make(chan bool, len(transports))
	// For each transports, concurencialy update infos
	for _, t := range transports {
		go func(tran models.ITransport) {
			err := make(chan error, 1)
			done := make(chan bool, 1)
			go func() {
				// Recover from potential panic in agencies
				defer func() {
					if r := recover(); r != nil {
						err <- fmt.Errorf("Panic updating transport (%v) infos\n	==> %v", tran, r)
					}
				}()
				err := tran.UpdateInfo()
				if err != nil {
					errors <- err
				} else {
					done <- true
				}
			}()

			select {
			case d := <-done:
				gather <- d
			case e := <-err:
				errors <- e
			case <-time.After(3 * time.Second):
				errors <- fmt.Errorf("Time out after 3s on %v", tran)
			}
		}(t)
	}
	// Wait for each update to finish
	i := 0
	errs := make([]error, 0, len(transports))
	for range transports {
		select {
		case <-gather:
			continue
		case e := <-errors:
			errs = append(errs, e)
		}
		i++
		if i == len(transports) {
			break
		}
	}
	return errs
}
