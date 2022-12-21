package index

import (
	"crawler/pkg/router"
	"net/http"
)

// GetIndex Function to Show API Information
func GetIndex(w http.ResponseWriter, r *http.Request) {
	router.ResponseSuccess(w, "200", "Go Framework is running")
}

// GetHealth Function to Show Health Check Status
func GetHealth(w http.ResponseWriter, r *http.Request) {
	router.HealthCheck(w)
}
