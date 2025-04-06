package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "weather_forecast/gen"
)

type (
	AverageTemperatureServer struct {
		pb.WeatherServiceServer
		avgTempService averageTempService
	}

	averageTempService interface {
		GetAverage(ctx context.Context, cityName string) (float32, error)
	}
)

func (s AverageTemperatureServer) GetAverageWeather(ctx context.Context, req *pb.AverageWeatherRequest) (*pb.AverageWeatherResponse, error) {
	avgTemp, err := s.avgTempService.GetAverage(ctx, req.City)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get average weather: %v", err)
	}
	return &pb.AverageWeatherResponse{
		AverageTemperature: avgTemp,
	}, nil
}

func NewAverageTemperatureServer(avgS averageTempService) *AverageTemperatureServer {
	return &AverageTemperatureServer{
		avgTempService: avgS,
	}
}
