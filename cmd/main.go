package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	weatherservice "weather_forecast/gen"
	"weather_forecast/internal/clients"
	grpcserver "weather_forecast/internal/controller/grpc"
	httpserver "weather_forecast/internal/controller/http"
	"weather_forecast/internal/db"
	"weather_forecast/internal/repository"
	"weather_forecast/internal/router"
	"weather_forecast/internal/services"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := db.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewWeatherRepository(conn)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGrpcServer(ctx, repo)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startHttpServer(ctx, repo)
	}()

	wg.Wait()
}

func startGrpcServer(ctx context.Context, repo *repository.WeatherRepositoryImpl) {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	avgTempService := services.NewAverageTemperatureService(repo)

	weatherservice.RegisterWeatherServiceServer(s, grpcserver.NewAverageTemperatureServer(avgTempService))

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("gRPC server listening on port 8082")
	<-ctx.Done()
	log.Printf("gRPC server shutting down")
	s.GracefulStop()
	log.Printf("gRPC server stopped")
}

func startHttpServer(ctx context.Context, repo *repository.WeatherRepositoryImpl) {

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	weatherApiClient := clients.NewWeatherApiClient(&httpClient)

	currWeatherService := services.NewCurrentWeatherService(weatherApiClient, repo)

	weatherController := httpserver.NewWeatherController(currWeatherService)

	handler := router.NewRouter(weatherController)
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("HTTP server listening on port 8080")
	<-ctx.Done()
	log.Printf("HTTP server shutting down")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Printf("HTTP server stopped")
}
