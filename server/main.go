package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Request coming from the React app
type reqBody struct {
	Address string `json:"address"`
}

// What we send back to the React app
type respBody struct {
	Entered      string `json:"entered"`
	Standardized string `json:"standardized"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// POST /api/standardize
	// Reads { "address": "..." } and returns entered + standardized (via Google)
	r.Post("/api/standardize", func(w http.ResponseWriter, r *http.Request) {
		var in reqBody
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Address == "" {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		if apiKey == "" {
			http.Error(w, "server missing GOOGLE_MAPS_API_KEY", http.StatusInternalServerError)
			return
		}

		// Minimal Google request: treat the whole input as a single address line
		payload := map[string]any{
			"address": map[string]any{
				"addressLines": []string{in.Address},
			},
		}
		body, _ := json.Marshal(payload)

		url := "https://addressvalidation.googleapis.com/v1:validateAddress?key=" + apiKey
		res, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			http.Error(w, "validation failed", http.StatusBadGateway)
			return
		}

		var gresp map[string]any
		if err := json.NewDecoder(res.Body).Decode(&gresp); err != nil {
			http.Error(w, "bad response", http.StatusBadGateway)
			return
		}

		// Default to the entered text; overwrite if Google returns formattedAddress
		standardized := in.Address
		if result, ok := gresp["result"].(map[string]any); ok {
			if addr, ok := result["address"].(map[string]any); ok {
				if formatted, ok := addr["formattedAddress"].(string); ok && formatted != "" {
					standardized = formatted
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(respBody{
			Entered:      in.Address,
			Standardized: standardized,
		})
	})

	log.Println("API on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
