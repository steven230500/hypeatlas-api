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
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	_ "github.com/steven230500/hypeatlas-api/docs"

	relaysvc "github.com/steven230500/hypeatlas-api/modules/relay/domain/service"
	relayhttp "github.com/steven230500/hypeatlas-api/modules/relay/infra/http"
	relayrepo "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository"

	signalhttp "github.com/steven230500/hypeatlas-api/modules/signal/infra/http"
	signalrepo "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository"

	sharedgorm "github.com/steven230500/hypeatlas-api/shared/db"
	sharedhttp "github.com/steven230500/hypeatlas-api/shared/http"
	"github.com/steven230500/hypeatlas-api/shared/logger"
)

func main() {
	_ = godotenv.Load()
	log := logger.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	storage := os.Getenv("STORAGE") // "memory" | "postgres"

	// DB (GORM)
	var gdb *gorm.DB
	if storage == "postgres" {
		gdb = sharedgorm.Connect()
		sharedgorm.Migrate(gdb)
		_ = sharedgorm.RunSeeds(gdb) // opcional (datos demo)
		log.Info().Msg("gorm postgres initialized")
	}

	// RELAY
	if gdb == nil {
		log.Fatal().Msg("relay requires postgres database")
	}
	relayRepo := relayrepo.New(gdb)
	log.Info().Msg("relay repository: postgres")
	relayService := relaysvc.New(relayRepo)
	relayHandler := relayhttp.New(relayService)
	hypeMapHandler := relayhttp.NewHypeMapHandler(relayService)

	// SIGNAL
	if gdb == nil {
		log.Fatal().Msg("signal requires postgres database")
	}
	signalRepo := signalrepo.New(gdb)
	log.Info().Msg("signal repository: postgres")
	signalRouter := signalhttp.NewRouter(signalRepo)

	// Router raíz
	r := sharedhttp.NewRouter()
	r.Use(middleware.Logger) // log de cada request

	// 404 en JSON (para que jq no muera con "404" texto)
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not_found","where":"root"}`))
	})

	// Health raíz
	// @Summary Healthcheck
	// @Success 200 {string} string "ok"
	// @Router /healthz [get]
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	// -- verificación Riot Games --
	// @Summary Riot Games domain verification
	// @Description Domain verification file for Riot Games API production key
	// @Produce plain
	// @Success 200 {string} string "verification code"
	// @Router /riot.txt [get]
	r.Get("/riot.txt", func(w http.ResponseWriter, r *http.Request) {
		v := os.Getenv("RIOT_SITE_VERIFICATION")
		if v == "" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(v))
	})

	// Swagger UI + alias /openapi.json
	mountSwagger(r)

	// API v1
	v1 := chi.NewRouter()
	relayHandler.Register(v1)
	hypeMapHandler.Register(v1)

	// ⬇️ Prefijo final: /v1/signal/...
	v1.Mount("/signal", signalRouter)

	// Health duplicado en el v1 para validar prefijo
	v1.Get("/signal/riot/_health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"where":"v1-direct"}`))
	})

	// Montar /v1 en el router raíz
	r.Mount("/v1", v1)

	// DUMP del router raíz (paths COMPLETOS)
	_ = chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		fmt.Println("[ROOT]", method, route)
		return nil
	})

	log.Info().Str("port", port).Msg("api up")
	_ = http.ListenAndServe(":"+port, r)
}
