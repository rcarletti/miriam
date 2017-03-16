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

	return w.Main.Temp, w.Weather[0].Description
}
