package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/chirouter"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/usecase"
)

func APIRouter(cfg *APIConfig) *chirouter.Wrapper {
	r := chirouter.NewWrapper(chi.NewRouter())
	r.Get("/healthz", Healthz)
	r.HandleFunc("/reset", cfg.Reset)

	r.Method(http.MethodGet, "/hello/{name}", nethttp.NewHandler(usecase.NewInteractor(HelloName)))
	return r
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called healthz")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type validateChirpBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := validateChirpBody{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(params.Body) > 140 {
		fmt.Printf("Chirp too long: %s", params.Body)
	}
	w.WriteHeader(200)
}

type HelloInput struct {
	Name string `path:"name" minLength:"3"`
	Age  int    `query:"age" min:"0"`
}

type HelloOutput struct {
	Message string `json:"message"`
}

func HelloName(ctx context.Context, input HelloInput, output *HelloOutput) error {
	output.Message = "Hello " + input.Name + ", you are " + fmt.Sprint(input.Age) + " years old."
	return nil
}
