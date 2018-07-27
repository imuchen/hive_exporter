package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func getTemperature(binary string) float64 {
	out, err := exec.Command(binary).Output()
	if err != nil {
		log.Println(err)
	}
	outString := string(out)
	log.Println(outString)
	tempString := strings.Split(outString, "\n")[0]
	log.Println(tempString)
	tempStringNumber := strings.Split(tempString, "=")[1]
	temperature, err := strconv.ParseFloat(tempStringNumber, 64)
	if err != nil {
		log.Println(err)
	}
	return temperature
}

func main() {
	binary := os.Getenv("SENSOR_READ_BINARY")
	if binary == "" {
		log.Fatal("Please supply path to binary with environment variable " +
			"SENSOR_READ_BINARY")
	}
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "hive",
			Name:      "temperature_deg_c",
			Help:      "Temperature inside the bee hive, in degree Celsius.",
		},
		func() float64 { return getTemperature(binary) },
	)); err != nil {
		log.Fatal(err)
	}
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "hive",
			Name:      "rel_humidity_percent",
			Help:      "Relative humidity in per cent inside the bee hive.",
		},
		func() float64 { return float64(50.6) },
	)); err != nil {
		log.Fatal(err)
	}
	http.Handle("/metrics", prometheus.Handler())
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
