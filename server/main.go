package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type reqBody struct {
	Address string `json:"address"`
}
type respBody struct {
	Entered      string `json:"entered"`
	Standardized string `json:"standardized"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Minimal stub: just trims and echos back (we'll replace with Google API later)
	r.Post("/api/standardize", func(w http.ResponseWriter, r *http.Request) {
		var in reqBody
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Address == "" {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		out := respBody{
			Entered:      in.Address,
			Standardized: in.Address, // placeholder for now
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	})

	log.Println("API on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
