package models

import "fmt"

// IAgency - Interface that every agencies need to implement
type IAgency interface {
	ContainsPosition(p *Position) bool
	TransportsNearPosition(p *Position, radius float64) ([]ITransport, error)
}

// AgencyProto - Can be embedded by custom Agencies structs
// The gives some properties and to methods
// Each Agencies structs still need to implement TransportsNearPosition
type AgencyProto struct {
	ID          string   `json:"ID"`          // ID of the region (Country.City.Agency)
	Name        string   `json:"name"`        // Displayed name of the Agency
	URL         string   `json:"url"`         // The URL to the agency's website/app...
	Git         string   `json:"git"`         // The URL to the git repo
	Center      Position `json:"center"`      // Center of the Agency
	Radius      float64  `json:"radius"`      // Radius of the Agency in meters
	Types       []int    `json:"types"`       // The type of transports handled by the agency
	TypesString []string `json:"typesString"` // Name for the type of transports
	IconsURL    []string `json:"iconsURL"`    // URL to the transports types icons
}

// String - Stringify an agency
// @return the agency's ID
func (a AgencyProto) String() string {
	return fmt.Sprintf("%v: %v o--------> %v m", a.Name, a.Center, a.Radius)
}

// ContainsPosition - Tell if a position is in the bounds of the agency
// @return the boolean, true if the position is in the agency's bounds
// TODO - Handle the case where the Agency is on the latitude or longitude 0
func (a *AgencyProto) ContainsPosition(p *Position) bool {
	return a.Center.DistanceFrom(p) <= a.Radius
}
