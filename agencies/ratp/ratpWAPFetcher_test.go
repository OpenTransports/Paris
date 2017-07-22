package ratp

import (
	"fmt"
	"os"
	"testing"

	htmlquery "github.com/antchfx/xquery/html"
)

func TestExtractTimesRER(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_rerA.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	passages, err := extractInfoRER(doc)
	if err != nil ||
		len(passages) != 3 ||
		passages[0].Direction != "Poissy (TETE)" ||
		len(passages[0].Times) != 2 ||
		passages[1].Direction != "St Germain en Laye (ZEMA)" ||
		len(passages[1].Times) != 3 ||
		passages[2].Direction != "Cergy-Le Haut (UPIR)" ||
		len(passages[2].Times) != 1 ||
		passages[0].Times[0] != "21:43" ||
		passages[1].Times[0] != "21:52" ||
		passages[2].Times[0] != "21:55" ||
		passages[1].Times[1] != "22:04" ||
		passages[0].Times[1] != "22:11" ||
		passages[1].Times[2] != "22:19" {
		for _, p := range passages {
			fmt.Println(p)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractTimesBus(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_bus.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	passages, err := extractInfo(doc)
	if err != nil ||
		len(passages) != 2 ||
		passages[0].Direction != "Champigny Camping Intern." ||
		len(passages[0].Times) != 2 ||
		passages[1].Direction != "Joinville-Le-Pont RER" ||
		len(passages[1].Times) != 2 ||
		passages[0].Times[0] != "6 mn" ||
		passages[0].Times[1] != "37 mn" ||
		passages[1].Times[0] != "13 mn" ||
		passages[1].Times[1] != "43 mn" {
		for _, p := range passages {
			fmt.Println(p)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractTimesMetro(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_metro.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	passages, err := extractInfo(doc)
	if err != nil ||
		len(passages) != 1 ||
		passages[0].Direction != "Porte de Clignancourt" ||
		len(passages[0].Times) != 4 ||
		passages[0].Times[0] != "6 mn" ||
		passages[0].Times[1] != "8 mn" ||
		passages[0].Times[2] != "13 mn" ||
		passages[0].Times[3] != "17 mn" {
		for _, p := range passages {
			fmt.Println(p)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractTimesTram(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_tram.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	passages, err := extractInfo(doc)
	if err != nil ||
		len(passages) != 1 ||
		passages[0].Direction != "Porte de Vincennes" ||
		len(passages[0].Times) != 2 ||
		passages[0].Times[0] != "1 mn" ||
		passages[0].Times[1] != "9 mn" {
		for _, p := range passages {
			fmt.Println(p)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractTimesNoct(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_noct.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	passages, err := extractInfo(doc)
	if err != nil ||
		len(passages) != 2 ||
		passages[0].Direction != "Mairie de St Ouen-Metro" ||
		passages[1].Direction != "Bourg-La-Reine RER" ||
		len(passages[0].Times) != 1 ||
		len(passages[1].Times) != 1 ||
		passages[0].Times[0] != "88 mn" ||
		passages[1].Times[0] != "49 mn" {
		for _, p := range passages {
			fmt.Println(p)
		}
		fmt.Println(err)
		t.Fail()
	}
}
