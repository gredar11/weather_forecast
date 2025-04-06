package models

import "time"

type CityWeatherModel struct {
	CityName    string
	GeneralInfo string
	Temperature float32
	Date        time.Time
}
