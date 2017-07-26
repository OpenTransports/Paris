// Package ratp - data.ratp.fr
package ratp

import (
	"encoding/json"

	"github.com/OpenTransports/Paris/helpers"
	"github.com/OpenTransports/Paris/models"
)

type (
	ratpAgency struct {
		models.AgencyProto
	}

	ratpTransport struct {
		models.TransportProto
	}
)

// Agency - Agency object
var Agency = &ratpAgency{
	models.AgencyProto{
		ID:     "FR.Paris.RATP",
		Name:   "RATP",
		URL:    "https://ratp.fr",
		Git:    "https://github.com/opentransports/paris",
		Radius: 20000, // 20 Km
		Center: models.Position{
			Latitude:  48.856,
			Longitude: 2.35,
		},
		Types:       []int{models.Tram, models.Metro, models.Bus, models.Rail},
		TypesString: []string{models.TramString, models.MetroString, models.BusString, models.RailString},
	},
}

// UpdateInfo for RATP Transport struct
func (t *ratpTransport) UpdateInfo() error {
	// Get next passages for the transport
	var err error
	t.Passages, err = GetNextPassages(t)
	return err
}

// TransportsNearPosition - Find transports contained in a circle
// @param p: the center of the circle
// @param radius: the radius of the circle in meters
// @return the transports list
func (a *ratpAgency) TransportsNearPosition(p *models.Position, radius float64) ([]models.ITransport, error) {
	// Get the transports list from the db
	rawTransports, err := helpers.GetTransports(a.ID)
	if err != nil {
		return nil, err
	}
	// transports, err := a.UnMarshall(rawTransports)
	var transports []ratpTransport
	err = json.Unmarshal(rawTransports, &transports)
	if err != nil {
		return nil, err
	}
	// Init array of filtered transports
	filteredTransports := make([]ratpTransport, 0, 200)
	// Loop trough agencies transports to find the one that are in the radius limits
	// Only return ONE Transport by line (the closest)
	i := 0
	for j, t := range transports {
		if t.DistanceFrom(p) < radius {
			added := false
			for k, ft := range filteredTransports {
				if ft.Line == t.Line {
					added = true
					if ft.DistanceFrom(p) > t.DistanceFrom(p) {
						filteredTransports[k] = transports[j]
						break
					}
				}
			}
			if !added {
				filteredTransports = append(filteredTransports, transports[j])
				i++
			}
		}
	}

	abstractedTransports := make([]models.ITransport, 0, i)
	for j := range filteredTransports {
		abstractedTransports = append(abstractedTransports, &filteredTransports[j])
	}
	// Return the transport slice from 0 to i
	// Indexes after i are just empty spots in the original array
	return abstractedTransports[:i], nil
}
