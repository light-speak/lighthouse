package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/manor"
	"github.com/light-speak/lighthouse/net"
	"github.com/light-speak/lighthouse/resolve"
	"github.com/light-speak/lighthouse/utils"
)

func StartService(resolver resolve.Resolve) {
	resolve.R = resolver
	if env.LighthouseConfig.App.Environment == env.Development {
		// Initialize progress bar
		fmt.Print("\033[?25l")       // Hide cursor
		defer fmt.Print("\033[?25h") // Show cursor when done
		utils.SmoothProgress(0, 20, "Initializing resolver", time.Millisecond*100, false)
	}

	err := graphql.LoadSchema()
	if err != nil {
		log.Error().Msgf("Failed to load schema: %v", err)
		return
	}

	if env.LighthouseConfig.App.Environment == env.Development {
		utils.SmoothProgress(20, 40, "Loading GraphQL schema", time.Millisecond*100, false)
	}

	port := env.LighthouseConfig.Server.Port

	if env.LighthouseConfig.App.Environment == env.Development {
		utils.SmoothProgress(40, 60, "Configuring server port", time.Millisecond*100, false)
	}

	r := net.New()

	if env.LighthouseConfig.App.Environment == env.Development {
		utils.SmoothProgress(60, 80, "Setting up router", time.Millisecond*100, false)
		utils.SmoothProgress(80, 100, "Starting server", time.Millisecond*100, false)

		// Clear the progress bar
		fmt.Print("\033[2K\r") // Clear current line and return cursor to start

		// Print final startup message
		fmt.Printf("\n\n🚀 GraphQL service started on port %s\n", port)
		fmt.Printf("📡 You can access the service at http://localhost:%s/query\n", port)
		fmt.Printf("🎨 You can access the studio at http://localhost:%s/studio\n", port)
		fmt.Printf("🔥 You can access the pprof at http://0.0.0.0:%s/debug/pprof\n\n", port)
	} else {
		log.Info().Msgf("GraphQL service started on port %s", port)
	}

	go manor.Manorable(graphql.GetParser().NodeStore)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r)
	if err != nil {
		log.Error().Msgf("Failed to start GraphQL service: %v", err)
	}
}
