package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

type Metric struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Alert struct {
	Severity string  `json:"severity"`
	Message  string  `json:"message"`
	Metric   string  `json:"metric"`
	Value    float64 `json:"value"`
}

var metrics []Metric
var alerts []Alert

func main() {
	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/alerts", handleAlerts)
	http.HandleFunc("/health", handleHealth)

	// Start metric collection
	go collectMetrics()

	log.Println("SRE Monitor starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func handleAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func collectMetrics() {
	for {
		// Simulate metric collection
		metrics = append(metrics, Metric{
			Name:      "cpu_usage",
			Value:     getRandomValue(0, 100),
			Timestamp: time.Now(),
		})

		metrics = append(metrics, Metric{
			Name:      "memory_usage",
			Value:     getRandomValue(0, 100),
			Timestamp: time.Now(),
		})

		// Check for alerts
		checkAlerts()

		// Keep only last 1000 metrics
		if len(metrics) > 1000 {
			metrics = metrics[len(metrics)-1000:]
		}

		time.Sleep(5 * time.Second)
	}
}

func checkAlerts() {
	for _, metric := range metrics {
		if metric.Name == "cpu_usage" && metric.Value > 80 {
			alerts = append(alerts, Alert{
				Severity: "warning",
				Message:  "High CPU usage detected",
				Metric:   metric.Name,
				Value:    metric.Value,
			})
		}
		if metric.Name == "memory_usage" && metric.Value > 90 {
			alerts = append(alerts, Alert{
				Severity: "critical",
				Message:  "Critical memory usage",
				Metric:   metric.Name,
				Value:    metric.Value,
			})
		}
	}

	// Keep only last 100 alerts
	if len(alerts) > 100 {
		alerts = alerts[len(alerts)-100:]
	}
}

func getRandomValue(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}
