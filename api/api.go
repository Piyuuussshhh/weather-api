package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Piyuuussshhh/weather-api/weather"
)

func Route() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "[ERROR] Please use the /weather?lat={}&long={} endpoint.", http.StatusNotFound)
	})

	http.HandleFunc("GET /weather", func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		long := r.URL.Query().Get("long")

		if lat == "" || long == "" {
			http.Error(w, "[ERROR] Invalid latitude longitude values", http.StatusBadRequest)
			return
		}

		weather, err := weather.GetWeather(r.Context(), lat, long)
		if err != nil {
			http.Error(w, fmt.Sprintf("[ERROR] %s\n", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		if err := json.NewEncoder(w).Encode(*weather); err != nil {
			http.Error(w, "[ERROR] Could not encode weather data to JSON", http.StatusInternalServerError)
			return
		}
	})

	return http.ListenAndServe(":8080", nil)
}