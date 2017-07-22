package models

import "fmt"

// Create constantes to identify transports types
const (
	Tram    = iota // Tram - open air public transportation rail network
	Metro          // Metro - underground public transportation rail network
	RER            // RER - regional public transportation rail network
	Bus            // Bus - open air public transportation network
	Unknown        // Unknown - unknown type of public transportation rail network
)

// Create constantes. Used to :
// 	- match transports types to strings
// 	- stringify transport
const (
	MetroString   = "metro"   // MetroString
	TramString    = "tram"    // TramString
	RERString     = "rer"     // RERString
	BusString     = "bus"     // BusString
	UnknownString = "unknown" // UnknownString
)

type (
	// ITransport - Interface that every transports need to implement
	ITransport interface {
		UpdateInfo() error
		DistanceFrom(p2 *Position) float64
	}

	// Group -
	// TODO - How is it used/usefull ?
	Group struct {
		ID    string
		Name  string
		Image string
	}

	// Passage - struct for public transports times
	Passage struct {
		Direction string   `json:"direction"` // Direction of the passage
		Times     []string `json:"times"`     // Time, is array of string to support non numeric values
	}

	// TransportProto - Can be embedded by custom Transports structs
	// Gives some usefull properties and two methods
	// Each Transports struct still need to implement UpdateInfo
	TransportProto struct {
		ID       string     `json:"ID"`       // ID of the Transport, should be specific to the Agency
		AgencyID string     `json:"agencyID"` // ID of the associated agency
		Type     int        `json:"type"`     // String identifing the kind of transport
		Image    string     `json:"image"`    // The image to display that represent the transport
		Name     string     `json:"name"`     // The name of the transport, doesn't have to be unique
		Group    string     `json:"group"`    // The group of the transport, aka the line
		Position Position   `json:"position"` // Position of the transport
		Passages []*Passage `json:"passages"` // Next passage for public transports
		Count    []int      `json:"count"`    // Count for free service bikes
	}
)

// DistanceFrom - Compute the distance between the transport point and a position
// @param p: the position
// @return the distance
func (t *TransportProto) DistanceFrom(p *Position) float64 {
	return t.Position.DistanceFrom(p)
}

// String - String representation of a transport
// Format: "<transport type> - <transport name>"
// Exemple: "metro - plaisance"
// @return the string representation
func (t TransportProto) String() string {
	switch t.Type {
	case Metro:
		return fmt.Sprintf("%v (%v): %v", MetroString, t.Group, t.Name)
	case Tram:
		return fmt.Sprintf("%v (%v): %v", TramString, t.Group, t.Name)
	case Bus:
		return fmt.Sprintf("%v (%v): %v", BusString, t.Group, t.Name)
	case RER:
		return fmt.Sprintf("%v (%v): %v", RERString, t.Group, t.Name)
	default:
		return fmt.Sprintf("(%v): %v", t.Group, t.Name)
	}
}
