package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ddo/go-fast"
)

func measureSpeed() ([]float64, error) {

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
			// fmt.Printf("%f\n", speed)
		}
	}()

	err = fastCom.Measure(urls, updates)
	if err != nil {
		return nil, err
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

func main() {
	dt := time.Now()
	data, err := measureSpeed()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf(`{"ts": "%s", "speed": %f}`, dt.Format("2006-01-02T15:04:05-0700"), avg(data))
}
