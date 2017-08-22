// Package ratp - data.ratp.fr
package ratp

import (
	"github.com/OpenTransports/Paris/helpers"
	"github.com/OpenTransports/lib-go/models"
)

// Transports - List of all transports of the RATP agency
var Transports []models.Transport

// Agency - Agency object
var Agency = models.Agency{
	ID:     "FR.Paris.RATP",
	Name:   "RATP",
	URL:    "https://ratp.fr",
	Radius: 20000, // 20 Km
	Center: models.Position{
		Latitude:  48.856,
		Longitude: 2.35,
	},
	Types: map[int]models.TransportTypeInfo{
		models.Tram: models.TransportTypeInfo{
			Icon: helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/L_T.png",
		},
		models.Metro: models.TransportTypeInfo{
			Icon: helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/L_M.png",
		},
		models.Bus: models.TransportTypeInfo{
			Icon: helpers.ServerURL + "/medias/logoBus.svg",
		},
		models.Rail: models.TransportTypeInfo{
			Name: "RER",
			Icon: helpers.ServerURL + "/medias/logoRER.png",
		},
	},
}

// UpdateTransportInfo for RATP Transport struct
func UpdateTransportInfo(transport models.Transport) error {
	// Get next passages for the transport
	var err error
	transport.Informations, err = GetNextPassages(transport)
	// Prevent return null instead of an empty array when transforming in json
	if transport.Informations == nil {
		transport.Informations = []models.Information{}
	}
	return err
}

// TransportsNearPosition - Find transports contained in a circle
// @param p: the center of the circle
// @param radius: the radius of the circle in meters
// @return the transports list
func TransportsNearPosition(position models.Position, radius float64) ([]models.Transport, error) {
	// Init array of filtered transports
	filteredTransports := make([]models.Transport, 0, 200)
	// Loop trough agencies transports to find the one that are in the radius limits
	i := 0
	for j, t := range Transports {
		if t.DistanceFrom(position) < radius {
			filteredTransports = append(filteredTransports, Transports[j])
			i++
		}
	}
	// Return the transport slice from 0 to i
	// Indexes after i are just empty spots in the original array
	return filteredTransports[:i], nil
}
