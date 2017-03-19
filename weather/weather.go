package weather

import (
	"log"

	owm "github.com/briandowns/openweathermap"
)

func GetWeather(city string) (termperature float64, description string) {
	w, err := owm.NewCurrent("C", "it")
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName(city)

	var currW string

	switch w.Weather[0].Icon {
	case "01d", "01n": //clear sky
		currW = "01"
	case "02d", "02n": //few clouds
		currW = "02"
	case "03d", "03n": //scattered clouds
		currW = "03"
	case "04d", "04n": //few clouds
		currW = "04"
	case "09d", "09n": //shower rain
		currW = "05"
	case "10d", "10n": // rain
		currW = "06"
	case "11d", "11n": //thunderstorm
		currW = "07"
	case "13d", "13n": //snow
		currW = "08"
	case "50d", "50n": //mist
		currW = "09"

	}

	return w.Main.Temp, currW
}
