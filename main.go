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

	// Serve web client
	app.Get("/", func(ctx context.Context) {
		ctx.ServeFile("./public/index.html", false)
	})
	assetHandler := app.StaticHandler("./public", false, false)
	app.SPA(assetHandler)

	// Build api
	apiRoute := app.Party("/api")
	// /api/transports?latitude=...&longitude=...
	apiRoute.Get("/transports", api.GetTransports)
	// /api/agencies?latitude=...&longitude=...
	apiRoute.Get("/agencies", api.GetAgencies)

	// Set listening port
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	app.Run(siris.Addr(":"+port), siris.WithCharset("UTF-8"))
}
