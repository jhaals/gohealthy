/*
Package gohealthy allow you to register health checks that can be hooked up to your load balancer such as haproxy

Example

 package main

 import (
 	"github.com/jhaals/gohealthy"
 	"time"
 )

 type customerPortal struct{}

 func (f customerPortal) GetHealth() gohealthy.HealthStatus {
 	return gohealthy.HealthStatus{Name: "customer-satisfaction",
 		Message: "Customer Satisfaction is above 180%",
 		Healthy: true}
 }

 type mySQLMonitor struct{}

 func (b mySQLMonitor) GetHealth() gohealthy.HealthStatus {
 	hs := gohealthy.HealthStatus{Name: "mysql-monitor"}
 	if time.Now().Weekday() == time.Friday {
 		hs.Message = "this is clearly not alright"
 		hs.Healthy = false
		return hs
 	}
 	hs.Message = "its not Friday so everything is awesome"
 	hs.Healthy = true
 	return hs
 }

 func main() {
 	g := gohealthy.Init()
 	g.Register(customerPortal{})
 	g.Register(mySQLMonitor{})
 	g.RunServer(1337)
 }

*/
package gohealthy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HealthStatus contains Name, Message and Health
type HealthStatus struct {
	Name    string `json:"-"`
	Message string `json:"message"`
	Healthy bool   `json:"healthy"`
}

// HealthCheck Status() need to be implemented by services that want to display health
type HealthCheck interface {
	GetHealth() HealthStatus
}

// HealthChecks contains list of HealthChecks
type HealthChecks struct {
	HealthChecks []HealthCheck
}

// Init a new instance of goheatlhy
func Init() HealthChecks {
	return HealthChecks{}
}

// Register new healthCheck
func (hcs *HealthChecks) Register(healthCheck HealthCheck) {
	hcs.HealthChecks = append(hcs.HealthChecks, healthCheck)
}

// RunServer starts a webserver serving health check status on specified port
func (hcs *HealthChecks) RunServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hcs.HealthCheckHandler(w, r)
	})
	log.Println(fmt.Sprintf("Serving health checks on port %v", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

/*
GetHealthChecks returns a map with all registered health checks and a bool
which will be false if one or more of the health checks return false
*/
func (hcs *HealthChecks) GetHealthChecks() (map[string]HealthStatus, bool) {
	var m = make(map[string]HealthStatus)
	heatlhy := true
	for _, hc := range hcs.HealthChecks {
		healthCheck := hc.GetHealth()
		if !healthCheck.Healthy {
			heatlhy = false
		}
		m[healthCheck.Name] = healthCheck
	}
	return m, heatlhy
}

/*
HealthCheckHandler returns a json with all registered healthChecks.
Response code will be set to 500 if one or more of the health checks are unhealthy
*/
func (hcs *HealthChecks) HealthCheckHandler(response http.ResponseWriter, request *http.Request) {
	healthChecks, healthy := hcs.GetHealthChecks()
	result, err := json.Marshal(healthChecks)
	if err != nil {
		http.Error(response, "Failed to json encode health status", 500)
	}
	if !healthy {
		http.Error(response, string(result), 500)
		return
	}
	response.Write(result)
}
