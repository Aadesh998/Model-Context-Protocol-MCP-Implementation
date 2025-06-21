package tools

import (
	"main/utils"
)

type GetWeatherInput struct {
	Location string `json:"location"`
}

var GetWeatherInputSchema = utils.GenerateSchema[GetWeatherInput]()

type GetWeatherResponse struct {
	Weather     string  `json:"weather"`
	Temperature float64 `json:"temperature"`
}

func GetWeather(location string) GetWeatherResponse {
	if location == "San Francisco" {
		return GetWeatherResponse{
			Weather:     "Sunny",
			Temperature: 21.0,
		}
	}
	return GetWeatherResponse{
		Weather:     "Unknown",
		Temperature: 0.0,
	}
}
