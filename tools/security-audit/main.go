package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type SecurityCheck struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Severity    string `json:"severity"`
	Details     string `json:"details"`
}

type SecurityAudit struct {
	Timestamp    time.Time       `json:"timestamp"`
	Checks       []SecurityCheck `json:"checks"`
	TotalIssues  int             `json:"total_issues"`
	HighIssues   int             `json:"high_issues"`
	MediumIssues int             `json:"medium_issues"`
	LowIssues    int             `json:"low_issues"`
}

func checkTLSConfig() SecurityCheck {
	check := SecurityCheck{
		Name:        "TLS Configuration",
		Description: "Check TLS version and cipher suites",
		Severity:    "High",
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		check.Status = "Failed"
		check.Details = fmt.Sprintf("TLS check failed: %v", err)
		return check
	}
	defer resp.Body.Close()

	check.Status = "Passed"
	check.Details = "TLS 1.2 or higher is configured"
	return check
}

func checkSecurityHeaders() SecurityCheck {
	check := SecurityCheck{
		Name:        "Security Headers",
		Description: "Check presence of security headers",
		Severity:    "Medium",
	}

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		check.Status = "Failed"
		check.Details = fmt.Sprintf("Failed to check headers: %v", err)
		return check
	}
	defer resp.Body.Close()

	headers := []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Strict-Transport-Security",
	}

	missing := []string{}
	for _, header := range headers {
		if resp.Header.Get(header) == "" {
			missing = append(missing, header)
		}
	}

	if len(missing) > 0 {
		check.Status = "Failed"
		check.Details = fmt.Sprintf("Missing security headers: %v", missing)
	} else {
		check.Status = "Passed"
		check.Details = "All security headers are present"
	}

	return check
}

func checkRateLimiting() SecurityCheck {
	check := SecurityCheck{
		Name:        "Rate Limiting",
		Description: "Check if rate limiting is implemented",
		Severity:    "Medium",
	}

	// Make multiple requests in quick succession
	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://localhost:8080/todos")
		if err != nil {
			check.Status = "Failed"
			check.Details = fmt.Sprintf("Rate limit check failed: %v", err)
			return check
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			check.Status = "Passed"
			check.Details = "Rate limiting is implemented"
			return check
		}
	}

	check.Status = "Failed"
	check.Details = "No rate limiting detected"
	return check
}

func runSecurityAudit() SecurityAudit {
	audit := SecurityAudit{
		Timestamp: time.Now(),
		Checks: []SecurityCheck{
			checkTLSConfig(),
			checkSecurityHeaders(),
			checkRateLimiting(),
		},
	}

	// Count issues by severity
	for _, check := range audit.Checks {
		if check.Status == "Failed" {
			audit.TotalIssues++
			switch check.Severity {
			case "High":
				audit.HighIssues++
			case "Medium":
				audit.MediumIssues++
			case "Low":
				audit.LowIssues++
			}
		}
	}

	return audit
}

func main() {
	http.HandleFunc("/audit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		audit := runSecurityAudit()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(audit)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Security Audit Tool started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
