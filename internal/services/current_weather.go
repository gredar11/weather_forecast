package services

import (
	"context"
	"sync"
	"time"
	"weather_forecast/internal/clients"
	. "weather_forecast/internal/models"
)

type CurrentWeatherService struct {
	currWeatherApi        currWeatherApi
	currWeatherRepository currWeatherRepository
	cache                 sync.Map
}

type weatherCacheKey struct {
	city string
	time int64
}

type (
	currWeatherApi interface {
		GetWeatherInCity(ctx context.Context, cityName string) (*clients.CityWeatherResp, error)
	}
	currWeatherRepository interface {
		AddCurrentWeather(ctx context.Context, cityName string, tmp float32) error
	}
)

func NewCurrentWeatherService(api currWeatherApi, repo currWeatherRepository) *CurrentWeatherService {
	return &CurrentWeatherService{
		currWeatherApi:        api,
		cache:                 sync.Map{},
		currWeatherRepository: repo,
	}
}

func (s *CurrentWeatherService) GetWeatherInCity(ctx context.Context, cityName string) (*CityWeatherModel, error) {
	if res, ok := s.getFromCache(cityName); ok {
		return res, nil
	}

	resp, err := s.currWeatherApi.GetWeatherInCity(ctx, cityName)
	if err != nil {
		return nil, err
	}

	cityWeatherModel := &CityWeatherModel{
		CityName:    resp.Location.Name,
		GeneralInfo: resp.RealTimeInfo.GeneralInfo.Text,
		Temperature: resp.RealTimeInfo.Temperature,
		Date:        resp.RealTimeInfo.MeasurementTime.Time,
	}

	if err = s.currWeatherRepository.AddCurrentWeather(ctx, cityWeatherModel.CityName, cityWeatherModel.Temperature); err != nil {
		return nil, err
	}

	s.addToCache(*cityWeatherModel)

	return cityWeatherModel, nil
}

func (s *CurrentWeatherService) getFromCache(city string) (*CityWeatherModel, bool) {
	key := weatherCacheKey{
		city: city,
		time: roundTimeToHour(time.Now().Local()),
	}
	val, ok := s.cache.Load(key)
	if !ok {
		return nil, false
	}
	if model, ok := val.(CityWeatherModel); ok {
		return &model, true
	}
	return nil, false
}

func (s *CurrentWeatherService) addToCache(model CityWeatherModel) {
	modelDate := model.Date.UTC()
	key := weatherCacheKey{
		city: model.CityName,
		time: roundTimeToHour(modelDate),
	}
	s.cache.Store(key, model)
}

func roundTimeToHour(t time.Time) int64 {
	hourRoundedTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	return hourRoundedTime.Unix()
}
