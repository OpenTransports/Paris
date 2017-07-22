package models

import "testing"

func createExampleAgency() *AgencyProto {
	// Create an agency
	return &AgencyProto{
		ID:   "Country.City.AgencyName",
		Name: "AgencyName",
		Center: Position{
			Latitude:  48.856,
			Longitude: 2.35,
		},
		Radius: 20000, // 20 Km
	}
}

func TestContainsPosition(t *testing.T) {
	ExampleAgency := createExampleAgency()
	// Try some points
	OK1 := ExampleAgency.ContainsPosition(&Position{48.856 /* GOOD */, 2.35 /* GOOD */})
	OK2 := ExampleAgency.ContainsPosition(&Position{48.8 /* GOOD */, 2.1 /* GOOD */})
	KO1 := ExampleAgency.ContainsPosition(&Position{49.3 /* GOOD */, 1.7 /* BAD  */})
	KO2 := ExampleAgency.ContainsPosition(&Position{48.2 /* BAD  */, 2.8 /* GOOD */})

	if !OK1 || !OK2 || KO1 || KO2 {
		t.Fail()
	}
}

func TestString(t *testing.T) {
	if createExampleAgency().String() != "AgencyName: (48.856 - 2.35) o--------> 20000 m" {
		t.Fail()
	}
}

// func TestTransportsNearPositionn(t *testing.T) {
// 	ExampleAgency := createExampleAgency()
// 	// Try some points
// 	OK1 := ExampleAgency.ContainsPosition(&Position{48.856 /* GOOD */, 2.35 /* GOOD */})
// 	OK2 := ExampleAgency.ContainsPosition(&Position{48.8 /* GOOD */, 2.1 /* GOOD */})
// 	KO1 := ExampleAgency.ContainsPosition(&Position{49.3 /* GOOD */, 1.7 /* BAD  */})
// 	KO2 := ExampleAgency.ContainsPosition(&Position{48.2 /* BAD  */, 2.8 /* GOOD */})
//
// 	if !OK1 || !OK2 || KO1 || KO2 {
// 		t.Fail()
// 	}
// }
