package ratp

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/OpenTransports/Paris/helpers"
	"github.com/OpenTransports/lib-go/models"
	"github.com/artonge/go-gtfs"
	"github.com/hashicorp/go-getter"
)

const gtfsURL = "http://dataratp.download.opendatasoft.com/RATP_GTFS_LINES.zip"
const iconFerreURL = "https://data.ratp.fr/api/datasets/1.0/indices-et-couleurs-de-lignes-du-reseau-ferre-ratp/attachments/indices_ferres_2017_05_zip.zip"
const iconBusURL = "http://data.ratp.fr/api/datasets/1.0/indices-des-lignes-de-bus-du-reseau-ratp/attachments/indices_zip.zip"
const iconLogoBusURL = "https://upload.wikimedia.org/wikipedia/commons/4/49/Paris_logo_bus_jms.svg"
const iconLogoRERURL = "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b0/Paris_RER_icon.svg/50px-Paris_RER_icon.svg.png"

var tmp = helpers.TmpDir(Agency.ID)
var media = helpers.MediaDir(Agency.ID)

func init() {
	// If testing don't do the setup
	if flag.Lookup("test.v") != nil {
		return
	}
	download(gtfsURL, tmp)
	unzip()
	load()
	go download(iconFerreURL, path.Join(media, "ferre"))
	go download(iconBusURL, path.Join(media, "bus"))
	go downloadFile(iconLogoRERURL, path.Join(media, "logoRER.png"))
	go downloadFile(iconLogoBusURL, path.Join(media, "logoBus.svg"))
}

// Download an file into the given path
// Skip if allready there
// go-getter handle the first unziping if needed
func download(URL string, path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Downloading in " + path + "...")
		err = getter.Get(path, URL)
		if err != nil {
			panic(err)
		}
	}
}
func downloadFile(URL string, path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Downloading in " + path + "...")
		err = getter.GetFile(path, URL)
		if err != nil {
			panic(err)
		}
	}
}

// Unzip the sub folders
func unzip() {
	fmt.Println("Unzipping RATP...")
	files, err := ioutil.ReadDir(tmp)
	if err != nil {
		panic(err)
	}
	for _, zip := range files {
		src := path.Join(tmp, zip.Name())
		dst := src[:len(src)-4]
		ext := src[len(src)-3:]
		if ext == "zip" {
			d := &getter.ZipDecompressor{}
			err = d.Decompress(dst, src, true)
			if err != nil {
				panic(err)
			}
			err = os.Remove(src)
			if err != nil {
				panic(err)
			}
		}
	}
}

func load() {
	fmt.Println("Loading RATP...")
	// Load data
	gtfss, err := gtfs.LoadSplitted(tmp, map[string]bool{"routes": true, "stops": true})
	if err != nil {
		panic(err)
	}
	Agency.Transports = removeDuplicate(mapToTransports(gtfss))
	Agency.Routes = mapToRoutes(gtfss, Agency.Transports)

}

// Remove Downloaded stuff
func clean() {
	fmt.Println("Cleaning RATP...")
	err := os.RemoveAll(tmp)
	if err != nil {
		panic(err)
	}
}

// Ratp data contains duplicates transports, this fontion filter them
func removeDuplicate(transports []models.Transport) []models.Transport {

	var isDuplicate bool
	var count int
	filteredTransports := make([]models.Transport, 0, len(transports))

	for _, transport := range transports {
		isDuplicate = false
		for _, fTransport := range filteredTransports {
			if transport.Line == fTransport.Line && transport.Name == fTransport.Name {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			filteredTransports = append(filteredTransports, transport)
			count++
		}
	}

	return filteredTransports[:count]
}

func mapToTransports(gtfss []*gtfs.GTFS) []models.Transport {
	// Total count of transports
	var size int
	for _, g := range gtfss {
		size += len(g.Stops)
	}
	// Create the transports array
	transports := make([]models.Transport, size)
	// For each gtfs, map the stops to a Transport struct
	// and add them to the transports array
	// Also update the image path of each transport depending on its Routes
	var i int
	for _, g := range gtfss {
		if g.Routes[0].ShortName == "T3" {
			g.Routes[0].ShortName = "T3a"
		}
		image := imageForRoute(g.Routes[0])
		for _, s := range g.Stops {
			transports[i] = models.Transport{
				ID:           s.ID,
				AgencyID:     Agency.ID,
				Name:         s.Name,
				Type:         g.Routes[0].Type,
				Line:         g.Routes[0].ShortName,
				IconURL:      image,
				Informations: []models.Information{},
				Position: models.Position{
					Latitude:  s.Latitude,
					Longitude: s.Longitude,
				},
			}
			i++
		}
	}
	return transports
}

// Map gtfs.Route to models.Route and fill models.Route.Stops with approriate transports
func mapToRoutes(gtfss []*gtfs.GTFS, transports []models.Transport) []models.Route {

	routes := make([]models.Route, 0, 200)
	for _, gtfs := range gtfss {
		for _, route := range gtfs.Routes {
			// Create the new route
			newRoute := models.Route{
				ID:     route.ShortName,
				Color:  route.Color,
				Points: make([]models.Position, 0, 50),
			}
			// Fill stops with related transports
			for _, transport := range transports {
				if transport.Line == newRoute.ID {
					newRoute.Points = append(newRoute.Points, transport.Position)
				}
			}
			// Append the new route to the routes array
			routes = append(routes, newRoute)
		}
	}

	return routes
}

// Given a route, return the corresponding path to the logo
func imageForRoute(r gtfs.Route) string {
	switch r.Type {
	case models.Tram:
		num, err := strconv.ParseInt(string(r.ShortName[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		if num >= 5 {
			return helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/t_" + r.ShortName[1:] + ".png"
		}
		return helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/T_" + r.ShortName[1:] + ".png"
	case models.Metro:
		return helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/M_" + r.ShortName + ".png"
	case models.Rail:
		return helpers.ServerURL + "/medias/ferre/indices-ferres-2017.05/RER_" + r.ShortName + ".png"
	case models.Bus:
		// Handle noctiliens Noct-133-genRVB
		if string(r.ShortName[0]) == "N" {
			return helpers.ServerURL + "/medias/bus/indices/Noct-" + r.ShortName[1:] + "-genRVB.png"
		}
		return helpers.ServerURL + "/medias/bus/indices/" + r.ShortName + "genRVB.png"
	default:
		return ""
	}
}
