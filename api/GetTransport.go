package api

import (
	"github.com/OpenTransports/Paris/agencies"
	"github.com/OpenTransports/lib-go/models"
	"github.com/go-siris/siris/context"
)

// GetTransport - /api/transports/transportID
// Get the freshly updated transport
// @param transportID : the id of a transports
func GetTransport(ctx context.Context) {
	// Recover from potential panic in agencies
	defer func() {
		if r := recover(); r != nil {
			ctx.Application().Logger().Errorf("Panic asking for nearest transports to agencies\n	==> %v", r)
			_, err := ctx.JSON(make([]models.Transport, 0))
			if err != nil {
				ctx.Application().Logger().Errorf("Error writting answer in /api/transports/transportID\n	==> %v", err)
			}
		}
	}()
	// Look for the transport in all agencies
	for _, agency := range agencies.AllAgencies {
		transport := agency.TransportForID(ctx.Params().Get("transportID"))
		if transport == nil {
			continue
		}
		// Get fresh infos for the transport
		errs := UpdateInfos([]models.Transport{*transport})
		for _, err := range errs {
			ctx.Application().Logger().Errorf("Error in /api/transports/transportID\n	==> %v", err)
		}
		// Return the result
		_, err := ctx.JSON(transport)
		// Log the error if any
		if err != nil {
			ctx.Application().Logger().Errorf("Error writting answer in /api/transports/transportID\n	==> %v", err)
		}
		break
	}
}
