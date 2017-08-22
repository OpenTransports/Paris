package agencies

import (
	"github.com/OpenTransports/Paris/agencies/ratp"
	"github.com/OpenTransports/lib-go/models"
)

// REGIONS - List of all regions
var agenciesList = [1]models.Agency{
	ratp.Agency,
}

// Containing -
func Containing(position models.Position) []models.Agency {

	if position.Latitude == 0 || position.Longitude == 0 {
		return agenciesList[:]
	}

	var filter = make([]models.Agency, 0)

	for _, a := range agenciesList {
		if a.ContainsPosition(position) {
			filter = append(filter, a)
		}
	}

	return filter
}
