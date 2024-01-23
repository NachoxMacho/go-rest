package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"nachoxmacho/go-rest/api"
	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/jsonschema"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/openapi"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/rest/response"
	"github.com/swaggest/swgui/v3cdn"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func fsServerRouter(cfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()
	r.Use(cfg.MiddlewareMetricsInc)
	r.Handle("/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	return r
}

func adminRouter(cfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()
	r.Get("/metrics", cfg.Metrics)
	return r
}

func main() {
	cfg := api.APIConfig{FileServerHits: 0}

	apiSchema := &openapi.Collector{}
	apiSchema.Reflector().SpecEns().Info.Title = "Basic Example"
	apiSchema.Reflector().SpecEns().Info.WithDescription("This app showcases a trivial REST API.")
	apiSchema.Reflector().SpecEns().Info.Version = "v1.2.3"

	// Setup request decoder and validator.
	validatorFactory := jsonschema.NewFactory(apiSchema, apiSchema)
	decoderFactory := request.NewDecoderFactory()
	decoderFactory.ApplyDefaults = true
	decoderFactory.SetDecoderFunc(rest.ParamInPath, chirouter.PathToURLValues)

	r := chirouter.NewWrapper(chi.NewRouter())
	r.Use(middlewareCors, nethttp.OpenAPIMiddleware(apiSchema), response.EncoderMiddleware, request.DecoderMiddleware(decoderFactory), request.ValidatorMiddleware(validatorFactory), middleware.Recoverer)
	r.Mount("/admin", adminRouter(&cfg))
	r.Mount("/app", fsServerRouter(&cfg))
	r.Mount("/api", api.APIRouter(&cfg))
	r.Get("/docs/openapi.json", apiSchema.ServeHTTP)

	swaggerUIHandler := v3cdn.NewHandler(apiSchema.Reflector().Spec.Info.Title, "/docs/openapi.json", "/docs")
	r.Mount("/docs", swaggerUIHandler)

	log.Print("Starting Server on :3333")

	err := http.ListenAndServe(":3333", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
