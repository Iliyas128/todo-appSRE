package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	todos = []Todo{
		{ID: 1, Task: "Initial Task", Completed: false, CreatedAt: time.Now()},
	}

	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method, path and status",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.Info("=== NEW BUILD: SRE DEBUG ===")

	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	json.NewEncoder(w).Encode(todos)
	log.WithFields(log.Fields{
		"method":   r.Method,
		"path":     r.URL.Path,
		"duration": time.Since(start).Seconds(),
	}).Info("Observe duration")
	requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())

	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"status": 200,
	}).Info("GET /todos success")
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		requestCounter.WithLabelValues(r.Method, r.URL.Path, "400").Inc()
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": 400,
			"error":  err.Error(),
		}).Error("Invalid JSON")
		return
	}
	newTodo.ID = len(todos) + 1
	newTodo.CreatedAt = time.Now()
	todos = append(todos, newTodo)

	json.NewEncoder(w).Encode(newTodo)
	log.WithFields(log.Fields{
		"method":   r.Method,
		"path":     r.URL.Path,
		"duration": time.Since(start).Seconds(),
	}).Info("Observe duration")
	requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())

	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"status": 200,
	}).Info("POST /todos success")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		srw := &statusRecorder{ResponseWriter: w, status: 200}
		start := time.Now()

		switch r.Method {
		case http.MethodGet:
			getTodos(srw, r)
		case http.MethodPost:
			createTodo(srw, r)
		default:
			http.Error(srw, "Method not allowed", http.StatusMethodNotAllowed)
			srw.status = http.StatusMethodNotAllowed
			log.WithFields(log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": srw.status,
			}).Warn("Unsupported HTTP method")
		}

		requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		requestCounter.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(srw.status)).Inc()
	})

	http.Handle("/metrics", promhttp.Handler())

	log.WithField("port", port).Info("Server started")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
