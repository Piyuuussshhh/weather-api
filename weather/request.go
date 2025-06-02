package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Piyuuussshhh/weather-api/cache"
)

const base_url = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

func GetWeather(ctx context.Context, cache *cache.Cache, lat, long string) (*Weather, error) {
	var weather Weather

	// If city data is stored in cache, return that.
	data, err := cache.GetCachedWeatherData(ctx, lat, long)
	if err != nil {
		return nil, err
	}
	if data != "" {
		if err := json.Unmarshal([]byte(data), &weather); err != nil {
			return nil, err
		}
		// fmt.Println("Data returned from redis cache!")
		return &weather, nil
	}

	// Else make api call.
	req_url := base_url + lat + "," + long + fmt.Sprintf("?unitGroup=metric&key=%s", os.Getenv("WEATHER_API_KEY"))
	response, err := http.Get(req_url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(response_body, &weather); err != nil {
		return nil, err
	}

	// Store the result in cache.
	if err := cache.CacheWeatherData(ctx, lat, long, string(response_body)); err != nil {
		return nil, err
	}

	// Return the result.
	// fmt.Println("Data returned from API.")
	return &weather, nil
}

type Weather struct {
	Location    string            `json:"resolvedAddress"`
	Description string            `json:"description"`
	Conditions  CurrentConditions `json:"currentConditions"`
}

type CurrentConditions struct {
	Datetime   string  `json:"datetime"`
	Temp       float64 `json:"temp"`
	Feelslike  float64 `json:"feelslike"`
	Humidity   float64 `json:"humidity"`
	Dew        float64 `json:"dew"`
	Precip     float64     `json:"precip"`
	Precipprob float64     `json:"precipprob"`
	Snow       float64     `json:"snow"`
	Snowdepth  float64     `json:"snowdepth"`
	Windgust   float64 `json:"windgust"`
	Windspeed  float64 `json:"windspeed"`
	Winddir    float64 `json:"winddir"`
	Visibility float64 `json:"visibility"`
	Uvindex    float64     `json:"uvindex"`
	Conditions string  `json:"conditions"`
	Sunrise    string  `json:"sunrise"`
	Sunset     string  `json:"sunset"`
}
