package models

import (
	"fmt"
	"math"
)

// Position - Structure containing latitude and longitude coordinates
type Position struct {
	Latitude  float64 `json:"latitude" query:"latitude"`   // Nord-South position
	Longitude float64 `json:"longitude" query:"longitude"` // East-West position
}

// String - Stringify a Position
func (p Position) String() string {
	return fmt.Sprintf("(%v - %v)", p.Latitude, p.Longitude)
}

// DistanceFrom - Compute the distance between two positions in meters
// See http://www.movable-type.co.uk/scripts/latlong.html
// @param p: the position
// @return the distance in meters
func (p *Position) DistanceFrom(p2 *Position) float64 {
	// Get radian dif between latitude and longitude the position
	dLat := toRadians(p2.Latitude - p.Latitude)
	dLon := toRadians(p2.Longitude - p.Longitude)
	// Some complexe computations with sin and cos
	a := math.Sin(dLat / 2)
	b := math.Sin(dLon / 2)
	c := math.Cos(toRadians(p.Latitude)) * math.Cos(toRadians(p2.Latitude))
	d := (a * a) + c*(b*b)
	// 6371000 => Average earth radius in meters
	return 6371000 * 2 * math.Atan2(math.Sqrt(d), math.Sqrt(1-d))
}

// Convert a degres angle to radian angle
// Exemple: 180° ==> π rad
func toRadians(degre float64) float64 {
	return degre * math.Pi / 180
}
