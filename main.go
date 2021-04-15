package main

import (
	"fmt"
	"github.com/bradenn/monoverse/simulation"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func main() {
	start := time.Now()
	fmt.Printf("BIG BANG! [%s]\n", time.Since(start))
	go func() {
		http.Handle("/", promhttp.Handler())
		_ = http.ListenAndServe(":2112", nil)
	}()

	sim := simulation.NewSimulation()
	fmt.Printf("Created The Universe [%s]\n", time.Since(start))
	sim.Run()

}
