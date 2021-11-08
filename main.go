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
var outputFormat string
var writer Writer

func init() {
	flag.BoolVar(&quiet, "q", false, "no status updates are printed to stdout")
	flag.StringVar(&outputFormat, "output", "json", "output format: json, csv")
	flag.Parse()
}

func measureSpeed() ([]float64, error) {

	if !quiet {
		fmt.Fprintf(os.Stderr, "Initialise measurements...\n")
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
				fmt.Fprintf(os.Stderr, "%f\n", speed)
			}
		}
	}()

	err = fastCom.Measure(urls, updates)
	if err != nil {
		return nil, err
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "Done.\n")
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

type Writer interface {
	WriteMeasurement(time time.Time, speed float64)
}

type JSONWriter struct{}

func (w JSONWriter) WriteMeasurement(time time.Time, speed float64) {
	ts := time.Format("2006-01-02T15:04:05-0700")
	if math.IsNaN(speed) {
		fmt.Fprintf(os.Stdout, `{"ts": "%s", "speed": null}`, ts)
	} else {
		fmt.Fprintf(os.Stdout, `{"ts": "%s", "speed": %f}`, ts, speed)
	}
}

type CSVWriter struct{}

func (w CSVWriter) WriteMeasurement(time time.Time, speed float64) {
	ts := time.Format("2006-01-02T15:04:05-0700")
	if math.IsNaN(speed) {
		fmt.Fprintf(os.Stdout, "%s;\n", ts)
	} else {
		fmt.Fprintf(os.Stdout, "%s;%f\n", ts, speed)
	}
}

func main() {
	dt := time.Now()

	switch outputFormat {
	case "json":
		writer = JSONWriter{}
	case "csv":
		writer = CSVWriter{}
	default:
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Error: Unsupported format %v\n", outputFormat)
		os.Exit(1)
	}

	data, err := measureSpeed()
	if err != nil {
		// errors should also be registered
		writer.WriteMeasurement(dt, math.NaN())
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	writer.WriteMeasurement(dt, avg(data))
}
