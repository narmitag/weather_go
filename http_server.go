package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func TestHandler(dataPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var observations Observations
		var highest_temp float64
		var lowest_temp float64

		var highest_temp_year float64
		var lowest_temp_year float64

		low_line := make([]opts.LineData, 0)
		high_line := make([]opts.LineData, 0)
		var xaxis []string

		years := []int{2019, 2020, 2021, 2022}
		months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}

		for _, year := range years {

			lowest_temp_year = 99
			highest_temp_year = -99

			for _, month := range months {
				highest_temp = -99
				lowest_temp = 99

				err := filepath.Walk(dataPath+"/"+strconv.Itoa(year)+"/"+month,
					func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						//println(path)
						if filepath.Ext(path) == ".json" {

							observations = process_file(path)
							lowest_temp, highest_temp = find_temp(observations, lowest_temp, highest_temp)
							if lowest_temp != 99 {
								low_line = append(low_line, opts.LineData{Value: lowest_temp})
								high_line = append(high_line, opts.LineData{Value: highest_temp})
								xaxis = append(xaxis, fileNameWithoutExtension(filepath.Base(path)))
							}

						}
						return nil
					})
				if err != nil {
					log.Println(err)
				}

				if lowest_temp < lowest_temp_year {
					lowest_temp_year = lowest_temp
				}
				if highest_temp > highest_temp_year {
					highest_temp_year = highest_temp
				}
			}

		}
		line := charts.NewLine()
		// set some global options like Title/Legend/ToolTip or anything else
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
			charts.WithTitleOpts(opts.Title{
				Title: "High and Low Daily Temp",
			}))

		// Put data into instance
		line.SetXAxis(xaxis).
			AddSeries("Lowest Temp", low_line).
			AddSeries("Highest Temp", high_line).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
		line.Render(w)
	}
}
