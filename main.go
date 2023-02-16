package main

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	httpPtr := flag.Bool("http", false, "Start Http Server")
	dataPath := flag.String("datapath", "", "Path to the data files")

	flag.Parse()

	httpPort := os.Getenv("PORT")

	if *httpPtr {
		println("Listening for data on :" + httpPort)
		http.HandleFunc("/", DataHandler(*dataPath))
		http.HandleFunc("/health", alive)
		http.ListenAndServe(":"+httpPort, nil)
	} else {
		ExtractData(*dataPath)
	}
}
