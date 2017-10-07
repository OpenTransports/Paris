package api

import (
	"github.com/OpenTransports/Paris/agencies"
	"github.com/go-siris/siris/context"
)

// GetTransportRoute - /api/transports/<transportID>/route
// Return the route associated to the transport
// @param transportID : the id of a transports
func GetTransportRoute(ctx context.Context) {
	for _, agency := range agencies.AllAgencies {
		transport := agency.TransportForID(ctx.Params().Get("transportID"))
		if transport == nil {
			continue
		}
		route := agency.RouteForTransports(*transport)
		_, err := ctx.JSON(route)
		// Log the error if any
		if err != nil {
			ctx.Application().Logger().Errorf("Error writting answer in /api/transports/transportID/route\n	==> %v", err)
		}
		break
	}
}
