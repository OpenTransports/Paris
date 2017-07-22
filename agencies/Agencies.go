package agencies

import (
	"github.com/OpenTransports/Paris/agencies/ratp"
	"github.com/OpenTransports/Paris/models"
)

// REGIONS - List of all regions
var agenciesList = [1]models.IAgency{
	ratp.Agency,
}

// Containing -
func Containing(p *models.Position) []models.IAgency {

	if p.Latitude == 0 || p.Longitude == 0 {
		return agenciesList[:]
	}

	var filter = make([]models.IAgency, 0)

	for _, a := range agenciesList {
		if a.ContainsPosition(p) {
			filter = append(filter, a)
		}
	}

	return filter
}
