package models

import (
	"math"
	"testing"
)

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
