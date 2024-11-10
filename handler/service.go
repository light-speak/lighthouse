package handler

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/net"
	"github.com/light-speak/lighthouse/resolve"
)

func StartService(resolver resolve.Resolve) {
	resolve.R = resolver
	if env.LighthouseConfig.App.Environment == env.Development {
		// Initialize progress bar
		fmt.Print("\033[?25l")       // Hide cursor
		defer fmt.Print("\033[?25h") // Show cursor when done
		smoothProgress(0, 20, "Initializing resolver", time.Millisecond*100)
	}

	err := graphql.LoadSchema()
	if err != nil {
		log.Error().Msgf("Failed to load schema: %v", err)
		return
	}

	if env.LighthouseConfig.App.Environment == env.Development {
		smoothProgress(20, 40, "Loading GraphQL schema", time.Millisecond*100)
	}

	port := env.LighthouseConfig.Server.Port

	if env.LighthouseConfig.App.Environment == env.Development {
		smoothProgress(40, 60, "Configuring server port", time.Millisecond*100)
	}

	r := net.New()

	if env.LighthouseConfig.App.Environment == env.Development {
		smoothProgress(60, 80, "Setting up router", time.Millisecond*100)
		smoothProgress(80, 100, "Starting server", time.Millisecond*100)

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

	go Manorable(graphql.GetParser().NodeStore)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r)
	if err != nil {
		log.Error().Msgf("Failed to start GraphQL service: %v", err)
	}
}

func smoothProgress(start, end int, status string, duration time.Duration) {
	steps := 15 // Reduce steps for faster updates
	delay := duration / time.Duration(steps)
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinIdx := 0

	for i := 0; i <= steps; i++ {
		progress := start + (end-start)*i/steps

		// Update spinner index
		spinIdx = (spinIdx + 1) % len(spinChars)

		showProgress(progress, status, spinChars[spinIdx])
		time.Sleep(delay)
		fmt.Print("\r")
	}
}

func showProgress(percent int, status string, spinChar string) {
	width := 30

	// Calculate completed width
	completed := width * percent / 100

	// Clear the current line and move to start
	fmt.Print("\033[2K\r")

	// Print spinner and progress bar
	fmt.Printf("%s [", spinChar)
	fmt.Print("\033[36m") // Cyan color
	fmt.Print(strings.Repeat("█", completed))
	if completed < width {
		fmt.Print(strings.Repeat("░", width-completed))
	}
	fmt.Print("\033[0m") // Reset color
	fmt.Printf("] %3d%% %s", percent, status)
}
