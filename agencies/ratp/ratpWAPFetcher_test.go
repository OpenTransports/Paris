package ratp

import (
	"fmt"
	"os"
	"testing"

	htmlquery "github.com/antchfx/xquery/html"
)

func TestExtractContentRER(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_rerA.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	infos, err := extractInfoRER(doc)
	if err != nil ||
		len(infos) != 3 ||
		len(infos[0].Content) != 2 ||
		infos[0].Title != "Poissy (TETE)" ||
		infos[0].Content[0] != "21:43" ||
		infos[0].Content[1] != "22:11" ||
		len(infos[1].Content) != 3 ||
		infos[1].Title != "St Germain en Laye (ZEMA)" ||
		infos[1].Content[0] != "21:52" ||
		infos[1].Content[1] != "22:04" ||
		infos[1].Content[2] != "22:19" ||
		len(infos[2].Content) != 1 ||
		infos[2].Title != "Cergy-Le Haut (UPIR)" ||
		infos[2].Content[0] != "21:55" {
		for _, info := range infos {
			fmt.Println(info)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractContentBus(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_bus.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	infos, err := extractInfo(doc)
	if err != nil ||
		len(infos) != 2 ||
		len(infos[0].Content) != 2 ||
		infos[0].Title != "Champigny Camping Intern." ||
		infos[0].Content[0] != "6 mn" ||
		infos[0].Content[1] != "37 mn" ||
		len(infos[1].Content) != 2 ||
		infos[1].Title != "Joinville-Le-Pont RER" ||
		infos[1].Content[0] != "13 mn" ||
		infos[1].Content[1] != "43 mn" {
		for _, info := range infos {
			fmt.Println(info)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractContentMetro(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_metro.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	infos, err := extractInfo(doc)
	if err != nil ||
		len(infos) != 1 ||
		infos[0].Title != "Porte de Clignancourt" ||
		len(infos[0].Content) != 4 ||
		infos[0].Content[0] != "6 mn" ||
		infos[0].Content[1] != "8 mn" ||
		infos[0].Content[2] != "13 mn" ||
		infos[0].Content[3] != "17 mn" {
		for _, info := range infos {
			fmt.Println(info)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractContentTram(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_tram.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	infos, err := extractInfo(doc)
	if err != nil ||
		len(infos) != 1 ||
		len(infos[0].Content) != 2 ||
		infos[0].Title != "Porte de Vincennes" ||
		infos[0].Content[0] != "1 mn" ||
		infos[0].Content[1] != "9 mn" {
		for _, info := range infos {
			fmt.Println(info)
		}
		fmt.Println(err)
		t.Fail()
	}
}

func TestExtractContentNoct(t *testing.T) {
	file, _ := os.OpenFile("testFiles/regexp_test_noct.html", os.O_RDONLY, 777)
	doc, _ := htmlquery.Parse(file)
	infos, err := extractInfo(doc)
	if err != nil ||
		len(infos) != 2 ||
		len(infos[0].Content) != 1 ||
		infos[0].Title != "Mairie de St Ouen-Metro" ||
		infos[0].Content[0] != "88 mn" ||
		len(infos[1].Content) != 1 ||
		infos[1].Title != "Bourg-La-Reine RER" ||
		infos[1].Content[0] != "49 mn" {
		for _, info := range infos {
			fmt.Println(info)
		}
		fmt.Println(err)
		t.Fail()
	}
}
