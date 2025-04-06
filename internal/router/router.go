package router

import (
	"net/http"
	http2 "weather_forecast/internal/controller/http"
)

func NewRouter(weatherController *http2.WeatherController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /get-weather", weatherController.HandleGetByCityName)
	return mux
}
