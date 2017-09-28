// Package ratp - data.ratp.fr
package ratp

import (
	"time"

	"github.com/OpenTransports/Paris/helpers"
	"github.com/OpenTransports/lib-go/models"
)

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
	Types: map[int]models.TypeInfo{
		models.Tram: models.TypeInfo{
			Icon: helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/L_T.png",
		},
		models.Metro: models.TypeInfo{
			Icon: helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/L_M.png",
		},
		models.Bus: models.TypeInfo{
			Icon: helpers.ServerURL + "/medias/logoBus.svg",
		},
		models.Rail: models.TypeInfo{
			Name: "RER",
			Icon: helpers.ServerURL + "/medias/logoRER.png",
		},
	},
}

// UpdateTransportInfo for RATP Transport struct
func UpdateTransportInfo(transport *models.Transport) error {
	// Get next passages for the transport
	var err error
	transport.Informations, err = GetNextPassages(*transport)
	// Prevent return null instead of an empty array when transforming in json
	for i := range transport.Informations {
		transport.Informations[i].Timestamp = int(time.Now().Unix()) * 1000
	}
	return err
}
