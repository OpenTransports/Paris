package ratp

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/OpenTransports/Paris/models"
	htmlquery "github.com/antchfx/xquery/html"
)

// Given a Transports
// 1 - Fetch the corresponding wap.ratp.fr page
// 2 - Extract informations (done with the wquery package)
// 3 - Build []*models.Passage from the extracted info

// GetNextPassages -
func GetNextPassages(t *ratpTransport) (passages []*models.Passage, err error) {
	var URL string
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v\n	==> URL: %v", r, URL)
		}
	}()
	var infoA, infoR []*models.Passage
	var errA, errR error
	// Build reseau and line depending on the kind of transport
	var reseau, line string
	switch t.Type {
	case models.Tram:
		reseau = models.TramString
		line = fmt.Sprintf("B%s", t.Line)
	case models.Metro:
		reseau = "metro"
		line = fmt.Sprintf("M%s", t.Line)
	case models.Bus:
		if t.Line[0] == 'N' {
			reseau = "noctilien"
			line = fmt.Sprintf("B%s", t.Line)
		} else {
			reseau = "bus"
			line = fmt.Sprintf("B%s", t.Line)
		}
	case models.Rail:
		reseau = "rer"
		line = fmt.Sprintf("R%s", t.Line)
	}
	// Fetch html page
	URL = fmt.Sprintf(
		"http://wap.ratp.fr/siv/schedule?reseau=%v&lineid=%v&stationname=%v&directionsens=",
		url.QueryEscape(reseau),
		url.QueryEscape(line),
		url.QueryEscape(t.Name),
	)
	docA, err := htmlquery.LoadURL(URL + "A")
	if err != nil {
		return nil, fmt.Errorf("Error updating transport\n	==> %v\n	==> station %v\n	==> URL: %v", err, t, URL+"A")
	}
	var docR *html.Node
	if t.Type != models.Bus {
		docR, err = htmlquery.LoadURL(URL + "R")
		if err != nil {
			return nil, fmt.Errorf("Error updating transport\n	==> %v\n	==> station %v\n	==> URL: %v", err, t, URL+"R")
		}
	}
	// Extract infos from html page
	switch t.Type {
	case models.Rail:
		infoA, errA = extractInfoRER(docA)
		infoR, errR = extractInfoRER(docR)
	case models.Bus:
		infoA, errA = extractInfo(docA)
	case models.Tram, models.Metro:
		infoA, errA = extractInfo(docA)
		infoR, errR = extractInfo(docR)
	}
	if errA != nil {
		return nil, fmt.Errorf("Error updating transport\n	==> %v\n	==> station %v\n	==> URL: %v", errA, t, URL+"A")
	}
	if errR != nil {
		return nil, fmt.Errorf("Error updating transport\n	==> %v\n	==> station %v\n	==> URL: %v", errR, t, URL+"R")
	}
	return append(infoA, infoR...), nil
}

func extractInfoNoct(doc *html.Node) ([]*models.Passage, error) {
	// Extract html nodes
	directionsNodes := htmlquery.Find(doc, "html/body/div[@class='subtitle' and starts-with(text(), 'Direction')]/b/text()")
	timesNodes := htmlquery.Find(doc, "html/body/div[starts-with(@class, 'bg') and count(*)=1]/b/child::text()")
	// Extract string from html node
	var directions, times []string
	for _, dir := range directionsNodes {
		directions = append(directions, dir.Data)
		directions = append(directions, dir.Data)
	}
	for _, t := range timesNodes {
		times = append(times, t.Data)
	}
	// Merge directions and times
	return mergeDirTime(directions, times)
}

func extractInfo(doc *html.Node) ([]*models.Passage, error) {
	// Extract html nodes
	directionsNodes := htmlquery.Find(doc, "html/body/div[starts-with(@class, 'bg') and count(*) = 0]/text()")[1:]
	timesNodes := htmlquery.Find(doc, "html/body/div[starts-with(@class, 'schmsg')]/b/child::text()")
	// Extract string from html node
	var directions, times []string
	for _, dir := range directionsNodes {
		directions = append(directions, strings.TrimSpace(dir.Data)[3:])
	}
	for _, t := range timesNodes {
		times = append(times, t.Data)
	}
	// Merge directions and times
	return mergeDirTime(directions, times)
}

func extractInfoRER(doc *html.Node) ([]*models.Passage, error) {
	// Extract html nodes
	directionsNodes := htmlquery.Find(doc, "html/body/div[starts-with(@class, 'bg')]/child::text()")[1:]
	trainsNamesNodes := htmlquery.Find(doc, "html/body/div[starts-with(@class, 'schmsg')]/a/child::text()")
	timesNodes := htmlquery.Find(doc, "html/body/div[starts-with(@class, 'schmsg')]/b/child::text()")
	// Check all length are the same (must be)
	if len(directionsNodes) != len(trainsNamesNodes) {
		fmt.Println("Weird stuff RER")
		for _, d := range directionsNodes {
			fmt.Println(d.Data)
		}
		for _, t := range trainsNamesNodes {
			fmt.Println(t.Data)
		}
		return nil, fmt.Errorf("Weird stuff RER")
	}
	// Extract string from html node
	var directions, times []string
	for i, dir := range directionsNodes {
		directions = append(directions, fmt.Sprintf("%v (%v)", strings.TrimSpace(dir.Data)[3:], trainsNamesNodes[i].Data))
	}
	for _, t := range timesNodes {
		times = append(times, t.Data)
	}
	// Merge directions and times
	return mergeDirTime(directions, times)
}

func mergeDirTime(directions []string, times []string) ([]*models.Passage, error) {
	// Directions and times must be of the same length
	if len(directions) != len(times) {
		fmt.Println("Weird stuff")
		for _, d := range directions {
			fmt.Println(d)
		}
		for _, t := range times {
			fmt.Println(t)
		}
		return nil, fmt.Errorf("Weird stuff")
	}
	// Associate each direction with some times
	// [ dir1   : [ time1
	//   dir2   :   time2
	//   dir1 ] :   time3 ]
	// -----------
	// [ dir1 : [time1, time3]
	//   dir2 : [time2] ]
	passages := []*models.Passage{}
	for i, dir := range directions {
		// Check if a passage for the dir allready exist
		// If it exist, store the time in it
		storred := false
		for _, p := range passages {
			if p.Direction == dir {
				p.Times = append(p.Times, times[i])
				storred = true
				break
			}
		}
		// And continue to the next time
		if storred {
			continue
		}
		// Else create a new passage inited with the dir and time
		p := &models.Passage{
			Direction: dir,
			Times:     []string{times[i]},
		}
		passages = append(passages, p)
	}
	return passages, nil
}
