// Package ratp - data.ratp.fr
package ratp

import (
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

var Transports []ratpTransport

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
		IconsURL: []string{
			helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/L_T",
			helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/L_M",
			helpers.ServerURL + "/medias/logoRER.png",
			helpers.ServerURL + "/medias/logoBus.svg",
		},
	},
}

// UpdateInfo for RATP Transport struct
func (t *ratpTransport) UpdateInfo() error {
	// Get next passages for the transport
	var err error
	t.Passages, err = GetNextPassages(t)
	// Prevent return null instead of an empty array when transforming in json
	if t.Passages == nil {
		t.Passages = []*models.Passage{}
	}
	return err
}

// TransportsNearPosition - Find transports contained in a circle
// @param p: the center of the circle
// @param radius: the radius of the circle in meters
// @return the transports list
func (a *ratpAgency) TransportsNearPosition(p *models.Position, radius float64) ([]models.ITransport, error) {
	// Init array of filtered transports
	filteredTransports := make([]ratpTransport, 0, 200)
	// Loop trough agencies transports to find the one that are in the radius limits
	// Only return ONE Transport by line (the closest)
	i := 0
	for j, t := range Transports {
		if t.DistanceFrom(p) < radius {
			added := false
			for k, ft := range filteredTransports {
				if ft.Line == t.Line {
					added = true
					if ft.DistanceFrom(p) > t.DistanceFrom(p) {
						filteredTransports[k] = Transports[j]
						break
					}
				}
			}
			if !added {
				filteredTransports = append(filteredTransports, Transports[j])
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
