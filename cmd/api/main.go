//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.3 init --parseDependency --parseInternal -g cmd/api/main.go -o docs

// @title           HypeAtlas API
// @version         1.0
// @description     Relay (HypeMap), Signal (MetaLens) y Comps & Leagues.
// @schemes         https http
// @host            api.hypeatlas.app
// @BasePath        /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/steven230500/hypeatlas-api/docs" // paquete generado por swag

	// RELAY
	relayout "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
	relaysvc "github.com/steven230500/hypeatlas-api/modules/relay/domain/service"
	relayhttp "github.com/steven230500/hypeatlas-api/modules/relay/infra/http"
	relaymem "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/memory"
	relaypg "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/postgres"

	// SIGNAL
	signalout "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
	signalsvc "github.com/steven230500/hypeatlas-api/modules/signal/domain/service"
	signalhttp "github.com/steven230500/hypeatlas-api/modules/signal/infra/http"
	signalmem "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository/memory"
	signalpg "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository/postgres"

	sharedpg "github.com/steven230500/hypeatlas-api/shared/db"
	sharedhttp "github.com/steven230500/hypeatlas-api/shared/http"
	"github.com/steven230500/hypeatlas-api/shared/logger"
)

func main() {
	log := logger.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	storage := os.Getenv("STORAGE") // "memory" | "postgres"

	// DB opcional
	var pool *pgxpool.Pool
	if storage == "postgres" {
		pool = sharedpg.MustOpen()
		log.Info().Msg("postgres pool initialized")
	}

	// RELAY
	var relayRepo relayout.Repository
	if pool != nil {
		relayRepo = relaypg.New(pool)
		log.Info().Msg("relay repository: postgres")
	} else {
		relayRepo = relaymem.New()
		log.Info().Msg("relay repository: memory")
	}
	relayService := relaysvc.New(relayRepo)
	relayHandler := relayhttp.New(relayService)

	// SIGNAL
	var signalRepo signalout.Repository
	if pool != nil {
		signalRepo = signalpg.New(pool)
		log.Info().Msg("signal repository: postgres")
	} else {
		signalRepo = signalmem.New()
		log.Info().Msg("signal repository: memory")
	}
	signalService := signalsvc.New(signalRepo)
	signalHandler := signalhttp.New(signalService)

	// Router
	r := sharedhttp.NewRouter()

	// Health (documentado)
	// @Summary Healthcheck
	// @Success 200 {string} string "ok"
	// @Router /healthz [get]
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	// Swagger UI + alias /openapi.json
	mountSwagger(r)

	// API
	r.Route("/v1", func(v chi.Router) {
		relayHandler.Register(v)
		signalHandler.Register(v)
		if pool != nil {
			v.Route("/ingest", func(ix chi.Router) {
				ix.Use(sharedhttp.ApiKeyMiddleware)
				relayIngest := relayhttp.NewIngest(pool)
				signalIngest := signalhttp.NewIngest(pool)
				ix.Route("/relay", relayIngest.Register)
				ix.Route("/signal", signalIngest.Register)
			})
		}
	})

	log.Info().Str("port", port).Msg("api up")
	_ = http.ListenAndServe(":"+port, r)
}
