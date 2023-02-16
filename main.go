package main

import (
	"flag"
	"net/http"
)

func main() {
	httpPtr := flag.Bool("http", false, "Start Http Server")
	dataPath := flag.String("datapath", "", "Path to the data files")

	flag.Parse()

	if *httpPtr {
		println("Listening for data on :8081")
		http.HandleFunc("/", DataHandler(*dataPath))
		http.HandleFunc("/health", alive)
		http.ListenAndServe(":8081", nil)
	} else {
		ExtractData(*dataPath)
	}
}
