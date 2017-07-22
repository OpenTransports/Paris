package models

import (
	"math"
	"testing"
)

// Earth circumference in km
const (
	earthCircumference = 40075
	metersByDegree     = (earthCircumference / 360) * 1000
)

func TestDistanceFrom(t *testing.T) {
	p0 := Position{48.82, 2.33}
	p1 := Position{48.83, 2.33}
	p2 := Position{48.82, 2.34}

	dist := p0.DistanceFrom(&p1)
	if math.Abs(1-dist/1112) > 0.01 {
		t.Fail()
	}

	dist = p0.DistanceFrom(&p2)
	if math.Abs(1-dist/732) > 0.01 {
		t.Fail()
	}

	dist = p0.DistanceFrom(&p0)
	if dist != 0 {
		t.Fail()
	}
}

func TestPositionString(t *testing.T) {
	p := Position{46, 2}
	if p.String() != "(46 - 2)" {
		t.Fail()
	}
}
