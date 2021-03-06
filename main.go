package main

import (
	"os"

	"github.com/OpenTransports/Paris/api"
	"github.com/go-siris/siris"
	"github.com/go-siris/siris/context"
)

func main() {
	// Create new siris app
	app := siris.New()

	// Serve medias files
	app.StaticWeb("/medias", "./medias")

	app.Use(func(ctx context.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Next()
	})

	// Build api
	app.Get("/transports", api.GetTransports)
	app.Get("/transports/{transportID:string}", api.GetTransport)
	app.Get("/transports/{transportID:string}/route", api.GetTransportRoute)
	app.Get("/agencies", api.GetAgencies)
	app.Get("/routes", api.GetRoutes)

	// Set listening port
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	app.Run(siris.Addr(":"+port), siris.WithCharset("UTF-8"))
}
