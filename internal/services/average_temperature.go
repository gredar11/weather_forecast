package services

import (
	"context"
	"weather_forecast/internal/repository"
)

type (
	AverageTemperatureService struct {
		repo weatherRepository
	}

	weatherRepository interface {
		GetForWeek(ctx context.Context, cityName string) ([]repository.DayWeather, error)
	}
)

func NewAverageTemperatureService(repo weatherRepository) *AverageTemperatureService {
	return &AverageTemperatureService{repo: repo}
}

func (a *AverageTemperatureService) GetAverage(ctx context.Context, cityName string) (float32, error) {
	temps, err := a.repo.GetForWeek(ctx, cityName)
	if err != nil {
		return 0, err
	}
	tempByDay := map[int64][]float32{}
	for _, temp := range temps {
		dayUnix := roundTimeToHour(temp.Date)
		if val, ok := tempByDay[dayUnix]; ok {
			tempByDay[dayUnix] = append(val, temp.Temperature)
			continue
		}
		tempByDay[dayUnix] = []float32{temp.Temperature}
	}

	days := len(tempByDay)
	var totalSum float32 = 0.0
	for _, temp := range tempByDay {
		totalSum += sum(temp)
	}
	return totalSum / float32(days), nil
}

func sum(vals []float32) float32 {
	var sum float32 = 0

	for _, val := range vals {
		sum += val
	}
	return sum
}
