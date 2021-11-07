package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/ddo/go-fast"
)

var quiet bool

func init() {
	flag.BoolVar(&quiet, "q", false, "no status updates are printed to stdout")
	flag.Parse()
}

func measureSpeed() ([]float64, error) {

	if !quiet {
		fmt.Println("Initialise measurements...")
	}

	fastCom := fast.New()
	err := fastCom.Init()
	if err != nil {
		return nil, err
	}

	urls, err := fastCom.GetUrls()
	if err != nil {
		return nil, err
	}

	measurements := make([]float64, 0)
	updates := make(chan float64)

	go func() {
		for speed := range updates {
			measurements = append(measurements, speed)
			if !quiet {
				fmt.Printf("%f\n", speed)
			}
		}
	}()

	err = fastCom.Measure(urls, updates)
	if err != nil {
		return nil, err
	}

	if !quiet {
		fmt.Println("Done.")
	}

	return measurements, nil
}

func avg(numbers []float64) float64 {
	sum := 0.0
	for i := 0; i < len(numbers); i++ {
		sum += numbers[i]
	}
	return sum / float64(len(numbers))
}

func writeMeasurement(time time.Time, speed float64) {
	ts := time.Format("2006-01-02T15:04:05-0700")
	if math.IsNaN(speed) {
		fmt.Printf(`{"ts": "%s", "speed": null}`, ts)
	} else {
		fmt.Printf(`{"ts": "%s", "speed": %f}`, ts, speed)
	}
}

func main() {
	dt := time.Now()

	data, err := measureSpeed()
	if err != nil {
		// errors should also be registered
		writeMeasurement(dt, math.NaN())
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	writeMeasurement(dt, avg(data))
}
