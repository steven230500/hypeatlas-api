package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	relaysvc "github.com/steven230500/hypeatlas-api/modules/relay/domain/service"
	relayhttp "github.com/steven230500/hypeatlas-api/modules/relay/infra/http"
	relaymem "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/memory"

	signalsvc "github.com/steven230500/hypeatlas-api/modules/signal/domain/service"
	signalhttp "github.com/steven230500/hypeatlas-api/modules/signal/infra/http"
	signalmem "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository/memory"

	sharedhttp "github.com/steven230500/hypeatlas-api/shared/http"
	"github.com/steven230500/hypeatlas-api/shared/logger"
)

func main() {
	log := logger.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// RELAY
	relayRepo := relaymem.New()
	relayService := relaysvc.New(relayRepo)
	relayHandler := relayhttp.New(relayService)

	// SIGNAL
	signalRepo := signalmem.New()
	signalService := signalsvc.New(signalRepo)
	signalHandler := signalhttp.New(signalService)

	r := sharedhttp.NewRouter()

	// Health
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })

	// Versioned API
	r.Route("/v1", func(v chi.Router) {
		relayHandler.Register(v)  // /v1/relay/...
		signalHandler.Register(v) // /v1/signal/...
	})

	log.Info().Str("port", port).Msg("api up")
	http.ListenAndServe(":"+port, r)
}
