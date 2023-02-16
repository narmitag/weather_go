package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Observations struct {
	Observations []Observation `json:"observations"`
}

type Observation struct {
	StationID          string  `json:"stationID"`
	Tz                 string  `json:"tz"`
	ObsTimeUtc         string  `json:"obsTimeUtc"`
	ObsTimeLocal       string  `json:"obsTimeLocal"`
	Epoch              int     `json:"epoch"`
	Lat                float64 `json:"lat"`
	Lon                float64 `json:"lon"`
	SolarRadiationHigh float64 `json:"solarRadiationHigh"`
	UvHigh             float64 `json:"uvHigh"`
	WinddirAvg         int     `json:"winddirAvg"`
	HumidityHigh       float64 `json:"humidityHigh"`
	HumidityLow        float64 `json:"humidityLow"`
	HumidityAvg        float64 `json:"humidityAvg"`
	QcStatus           float64 `json:"qcStatus"`
	Metric             Metric  `json:"metric"`
}

type Metric struct {
	TempHigh      float64 `json:"tempHigh"`
	TempLow       float64 `json:"tempLow"`
	TempAvg       float64 `json:"tempAvg"`
	WindspeedHigh float64 `json:"windspeedHigh"`
	WindspeedLow  float64 `json:"windspeedLow"`
	WindspeedAvg  float64 `json:"windspeedAvg"`
	WindgustHigh  float64 `json:"windgustHigh"`
	WindgustLow   float64 `json:"windgustLow"`
	WindgustAvg   float64 `json:"windgustAvg"`
	DewptHigh     float64 `json:"dewptHigh"`
	DewptLow      float64 `json:"dewptLow"`
	DewptAvg      float64 `json:"dewptAvg"`
	WindchillHigh float64 `json:"windchillHigh"`
	WindchillLow  float64 `json:"windchillLow"`
	WindchillAvg  float64 `json:"windchillAvg"`
	HeatindexHigh float64 `json:"heatindexHigh"`
	HeatindexLow  float64 `json:"heatindexLow"`
	HeatindexAvg  float64 `json:"heatindexAvg"`
	PressureMax   float64 `json:"pressureMax"`
	PressureMin   float64 `json:"pressureMin"`
	PressureTrend float64 `json:"pressureTrend"`
	PrecipRate    float64 `json:"precipRate"`
	PrecipTotal   float64 `json:"precipTotal"`
}

func ExtractData(dataPath string) {

	var observations Observations
	var highest_temp float64
	var lowest_temp float64

	var highest_temp_year float64
	var lowest_temp_year float64

	years := []int{2019, 2020, 2021, 2022}
	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	headers := []string{"Timestamp", "tempLow", "tempHigh", "WindspeedHigh", "WindspeedLow",
		"WindspeedAvg", "WindgustHigh", "WindgustLow", "WindgustAvg", "PrecipRate", "PrecipTotal"}
	csvFile, err := os.Create("data.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	w := csv.NewWriter(csvFile)
	defer w.Flush()
	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, year := range years {

		println("")
		println("Year : " + strconv.Itoa(year))
		println("----------------------------")
		lowest_temp_year = 99
		highest_temp_year = -99

		for _, month := range months {
			highest_temp = -99
			lowest_temp = 99

			println("")
			println("+ Month : " + month)
			//data/ISTGBUCH2/hourly
			err := filepath.Walk(dataPath+"/"+strconv.Itoa(year)+"/"+month,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					//fmt.Println(path)
					if filepath.Ext(path) == ".json" {

						observations = process_file(path)
						lowest_temp, highest_temp = find_data(observations, lowest_temp, highest_temp, *w)

					}
					return nil
				})
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("  Highest Temp = %f\n", highest_temp)
			fmt.Printf("  Lowest Temp  = %f\n", lowest_temp)

			if lowest_temp < lowest_temp_year {
				lowest_temp_year = lowest_temp
			}
			if highest_temp > highest_temp_year {
				highest_temp_year = highest_temp
			}
		}

		println("")
		println("++ Summary for " + strconv.Itoa(year))

		fmt.Printf("   Highest Temp = %f\n", highest_temp_year)
		fmt.Printf("   Lowest Temp  = %f\n", lowest_temp_year)

	}
	w.Flush()
	csvFile.Close()
}

func find_data(observations Observations,
	lowest float64,
	highest float64, w csv.Writer) (float64, float64) {

	for i := 0; i < len(observations.Observations); i++ {
		if observations.Observations[i].Metric.TempHigh > highest {
			highest = observations.Observations[i].Metric.TempHigh
		}
		if observations.Observations[i].Metric.TempLow < lowest {
			lowest = observations.Observations[i].Metric.TempLow
		}
		// row := []string{observations.Observations[i].ObsTimeUtc,
		// 	strconv.FormatFloat(observations.Observations[i].Metric.TempLow, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.TempHigh, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.WindspeedHigh, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.WindspeedLow, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.WindspeedAvg, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.WindgustHigh, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.WindgustLow, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.WindgustAvg, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.PrecipRate, 'E', -1, 32),
		// 	strconv.FormatFloat(observations.Observations[i].Metric.PrecipTotal, 'E', -1, 32),
		// }
		row := []string{observations.Observations[i].ObsTimeUtc,
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.TempLow),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.TempHigh),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.WindspeedHigh),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.WindspeedLow),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.WindspeedAvg),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.WindgustHigh),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.WindgustLow),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.WindgustAvg),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.PrecipRate),
			fmt.Sprintf("%.2f", observations.Observations[i].Metric.PrecipTotal),
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	return lowest, highest
}
