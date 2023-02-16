package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func process_file(filename string) Observations {
	jsonFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var observations Observations
	json.Unmarshal(byteValue, &observations)

	return observations
}

func find_temp(observations Observations,
	lowest float64,
	highest float64) (float64, float64) {

	for i := 0; i < len(observations.Observations); i++ {
		if observations.Observations[i].Metric.TempHigh > highest {
			highest = observations.Observations[i].Metric.TempHigh
		}
		if observations.Observations[i].Metric.TempLow < lowest {
			lowest = observations.Observations[i].Metric.TempLow
		}
	}
	return lowest, highest
}
