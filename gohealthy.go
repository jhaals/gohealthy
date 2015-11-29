package gohealthy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HealthStatus contains Name, Message and Health
type HealthStatus struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Healthy bool   `json:"healthy"`
}

// HealthCheck Status() need to be implemented by services that want to display health
type HealthCheck interface {
	Status() HealthStatus
}

// HealthChecks contains list of HealthChecks
type HealthChecks struct {
	HealthChecks []HealthCheck
}

// Register new healthCheck
func (hcs *HealthChecks) Register(healthCheck HealthCheck) {
	hcs.HealthChecks = append(hcs.HealthChecks, healthCheck)
}

// Status of all healthChecks
func (hcs *HealthChecks) Status() ([]HealthStatus, bool) {
	var x []HealthStatus
	heatlhy := true
	for _, hc := range hcs.HealthChecks {
		if !hc.Status().Healthy {
			heatlhy = false
		}
		x = append(x, hc.Status())
	}
	return x, heatlhy
}

// StatusHandler to be used in webserver
func (hcs *HealthChecks) StatusHandler(response http.ResponseWriter, request *http.Request) {
	status, healthy := hcs.Status()
	result, err := json.Marshal(status)
	if err != nil {
		http.Error(response, "Failed to json encode health status", 500)
	}
	if !healthy {
		http.Error(response, string(result), 500)
		return
	}
	response.Write(result)
}

// Init a new instance of goheatlhy
func Init() HealthChecks {
	return HealthChecks{}
}

// RunServer starts a webserver serving healthcheck status on specified port
func (hcs *HealthChecks) RunServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hcs.StatusHandler(w, r)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
