package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type DayWeather struct {
	Temperature float32   `db:"temperature"`
	Date        time.Time `db:"date"`
}

type WeatherRepositoryImpl struct {
	conn *pgxpool.Pool
}

var NotFoundWeatherErr = errors.New("not found weather")

func (r WeatherRepositoryImpl) AddCurrentWeather(ctx context.Context, cityName string, tmp float32) error {
	query := "INSERT INTO daily_temperature (temperature, date, city_name) values ($1, $2, $3);"
	_, err := r.conn.Exec(ctx, query, tmp, time.Now(), cityName)
	return err
}

func (r WeatherRepositoryImpl) GetForWeek(ctx context.Context, cityName string) ([]DayWeather, error) {
	query := `SELECT temperature, date FROM daily_temperature 
			WHERE date > now() - INTERVAL '1 WEEK' AND city_name = $1;`
	rows, err := r.conn.Query(ctx, query, cityName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	temps := make([]DayWeather, 0)

	for rows.Next() {
		var dayWeather DayWeather
		scanErr := rows.Scan(&dayWeather.Temperature, &dayWeather.Date)
		if scanErr != nil {
			return nil, scanErr
		}
		temps = append(temps, dayWeather)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return temps, nil
}

func NewWeatherRepository(conn *pgxpool.Pool) *WeatherRepositoryImpl {
	return &WeatherRepositoryImpl{conn: conn}
}
