package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
	"github.com/steven230500/hypeatlas-api/modules/signal/domain/service"
	"github.com/steven230500/hypeatlas-api/providers/riot"
)

// NewRouter crea un router con todos los handlers del módulo signal
func NewRouter(repo out.Repository) chi.Router {
	r := chi.NewRouter()

	// Servicio principal del módulo
	signalSvc := service.New(repo)

	// Cliente de Riot (condicional por env)
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	fmt.Println("Checking RIOT_API_KEY:", riotAPIKey != "")
	var riotSvc *riot.Service
	var metaGameSvc *service.MetaGameService
	if riotAPIKey != "" {
		fmt.Println("RIOT_API_KEY found, initializing Riot services...")
		riotSvc = riot.NewService(riotAPIKey, repo)
		metaGameSvc = service.NewMetaGameService(repo, riotSvc)
		fmt.Println("Riot services initialized successfully")
	} else {
		fmt.Println("RIOT_API_KEY not found in environment")
	}

	// Handlers principales (sin prefijos internos)
	signalHandler := New(signalSvc)
	signalHandler.Register(r)

	// Handler de Riot (si hay key)
	if riotSvc != nil && metaGameSvc != nil {
		fmt.Println("Registering Riot handler...")
		riotHandler := NewRiotHandler(riotSvc, signalSvc, metaGameSvc)
		riotHandler.Register(r)
		r.Get("/riot/_health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true}`))
		})
		fmt.Println("Riot handler registered successfully")
	} else {
		fmt.Println("Riot services not available - handler not registered")
	}

	// Dump del SUBROUTER (paths relativos)
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		fmt.Printf("[ROUTE] %s %s\n", method, route)
		return nil
	})

	return r
}
