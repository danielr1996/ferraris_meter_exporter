package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"net/http"
)

const PIN = rpi.GPIO2
const PORT = "2112"
const CHIP = "gpiochip0"

var rpmCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ferraris_meter_rpms",
	Help: "Total rpms of ferraris meter",
})

func HandleFallingEdge() (*gpiod.Chip, *gpiod.Line){
	c, err := gpiod.NewChip(CHIP)
	if err != nil {
		panic(err)
	}
	l, err := c.RequestLine(PIN, gpiod.WithBothEdges(func (evt gpiod.LineEvent) {
		if evt.Type == gpiod.LineEventFallingEdge {
			rpmCounter.Inc()
			fmt.Println("Detected falling edge")
		}
	}))
	if err != nil {
		panic(err)
	}
	return c, l
}

func main() {
	c, l := HandleFallingEdge()
	defer c.Close()

	defer l.Close()

	fmt.Println("Starting ferraris_counter_exporter on port "+PORT)
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":"+PORT, nil)
}
