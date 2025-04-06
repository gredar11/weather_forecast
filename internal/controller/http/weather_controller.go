package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	. "weather_forecast/internal/models"
	"weather_forecast/internal/repository"
)

type (
	currentWeatherInCityService interface {
		GetWeatherInCity(ctx context.Context, cityName string) (*CityWeatherModel, error)
	}
)

type WeatherController struct {
	service currentWeatherInCityService
}

func NewWeatherController(service currentWeatherInCityService) *WeatherController {
	return &WeatherController{service}
}

func (c WeatherController) HandleGetByCityName(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("cityName")
	if cityName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currentWeather, err := c.service.GetWeatherInCity(r.Context(), cityName)
	if err != nil {
		if errors.Is(err, repository.NotFoundWeatherErr) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(currentWeather)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
