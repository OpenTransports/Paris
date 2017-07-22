package models

import (
	"math"
	"testing"
)

func TestTransportString(t *testing.T) {
	t1 := TransportProto{ID: "1", Type: Tram, Group: "11", Name: "Tram Name"}
	t2 := TransportProto{ID: "2", Type: Metro, Group: "22", Name: "Metro Name"}
	t3 := TransportProto{ID: "3", Type: Bus, Group: "33", Name: "Bus Name"}
	t4 := TransportProto{ID: "4", Type: Unknown, Group: "44", Name: "Unknown Name"}

	if t1.String() != "tram (11): Tram Name" ||
		t2.String() != "metro (22): Metro Name" ||
		t3.String() != "bus (33): Bus Name" ||
		t4.String() != "(44): Unknown Name" {
		t.Fail()
	}
}

func TestTransportDistanceFrom(t *testing.T) {
	tr := TransportProto{Position: Position{0, 0}}
	p := Position{0, 1}
	// Compute distance between the two points
	dist := tr.DistanceFrom(&p)
	// Test that the error is less than 1%
	err := math.Abs(1 - dist/metersByDegree)
	if err > 0.01 {
		t.Fail()
	}
}
