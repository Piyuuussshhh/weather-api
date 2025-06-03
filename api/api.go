package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Piyuuussshhh/weather-api/cache"
	"github.com/Piyuuussshhh/weather-api/middleware"
	"github.com/Piyuuussshhh/weather-api/weather"
	"golang.org/x/time/rate"
)

func Route(ctx context.Context) error {
	cache, err := cache.NewCache()
	if err != nil {
		return err
	}
	
	go func() {
		<-ctx.Done()
		if err := cache.Client.Close(); err != nil {
			fmt.Printf("[ERROR] Failed to close cache client: %v\n", err)
		} else {
			fmt.Println("[INFO] Cache client closed successfully.")
		}
	}()
	
	mux := http.NewServeMux()

	rl := middleware.NewRateLimiter(rate.Every(1*time.Second), 5) // 5 requests per second

	mux.HandleFunc("GET /weather", rl.Limit(
		func(w http.ResponseWriter, r *http.Request) {
			lat := r.URL.Query().Get("lat")
			long := r.URL.Query().Get("long")

			if lat == "" || long == "" {
				http.Error(w, "[ERROR] Invalid latitude longitude values", http.StatusBadRequest)
				return
			}

			weather, err := weather.GetWeather(r.Context(), cache, lat, long)
			if err != nil {
				http.Error(w, fmt.Sprintf("[ERROR] %s\n", err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(*weather); err != nil {
				http.Error(w, "[ERROR] Could not encode weather data to JSON", http.StatusInternalServerError)
				return
			}
		},
	))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("[ERROR] Failed to shutdown server: %v\n", err)
		} else {
			fmt.Println("[INFO] Server shutdown successfully.")
		}
	}()

	return server.ListenAndServe()
}