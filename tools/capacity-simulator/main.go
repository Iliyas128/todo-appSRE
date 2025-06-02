package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

type SimulationConfig struct {
	Duration     int     `json:"duration"`      // Duration in minutes
	UsersPerSec  float64 `json:"users_per_sec"` // Number of users per second
	GrowthRate   float64 `json:"growth_rate"`   // Monthly growth rate
	PeakFactor   float64 `json:"peak_factor"`   // Peak to average ratio
	ResponseTime int     `json:"response_time"` // Target response time in ms
}

type SimulationResult struct {
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	AvgResponseTime    float64 `json:"avg_response_time"`
	MaxResponseTime    float64 `json:"max_response_time"`
	MinResponseTime    float64 `json:"min_response_time"`
	P95ResponseTime    float64 `json:"p95_response_time"`
	P99ResponseTime    float64 `json:"p99_response_time"`
	CPUUsage           float64 `json:"cpu_usage"`
	MemoryUsage        float64 `json:"memory_usage"`
}

func simulateLoad(config SimulationConfig) SimulationResult {
	var result SimulationResult
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(config.Duration) * time.Minute)

	// Initialize metrics
	responseTimes := make([]float64, 0)

	// Simulate requests
	for time.Now().Before(endTime) {
		// Calculate current load with peak factor
		currentLoad := config.UsersPerSec * (1 + config.PeakFactor*math.Sin(float64(time.Now().Unix())/3600))

		// Simulate batch of requests
		for i := 0; i < int(currentLoad); i++ {
			// Simulate request
			start := time.Now()

			// Make actual request to the application
			resp, err := http.Get("http://localhost:8080/todos")
			if err != nil {
				result.FailedRequests++
				continue
			}
			resp.Body.Close()

			// Record metrics
			duration := float64(time.Since(start).Milliseconds())
			responseTimes = append(responseTimes, duration)

			if resp.StatusCode == http.StatusOK {
				result.SuccessfulRequests++
			} else {
				result.FailedRequests++
			}

			result.TotalRequests++
		}

		// Simulate resource usage
		result.CPUUsage = 30 + rand.Float64()*40      // 30-70% CPU usage
		result.MemoryUsage = 200 + rand.Float64()*100 // 200-300MB memory usage

		time.Sleep(time.Second)
	}

	// Calculate statistics
	if len(responseTimes) > 0 {
		sort.Float64s(responseTimes)
		result.AvgResponseTime = average(responseTimes)
		result.MaxResponseTime = responseTimes[len(responseTimes)-1]
		result.MinResponseTime = responseTimes[0]
		result.P95ResponseTime = responseTimes[int(float64(len(responseTimes))*0.95)]
		result.P99ResponseTime = responseTimes[int(float64(len(responseTimes))*0.99)]
	}

	return result
}

func average(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func main() {
	http.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var config SimulationConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := simulateLoad(config)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Capacity Planning Simulator started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
