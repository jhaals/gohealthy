package gohealthy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mock struct{}

func (m mock) GetHealth() HealthStatus {
	return HealthStatus{"TestRegister", "This is just a test", true}
}
func TestRegister(t *testing.T) {

	goheatlhy := HealthChecks{}
	goheatlhy.Register(mock{})

	healthChecks, status := goheatlhy.GetHealthChecks()
	if !status {
		t.Errorf("this should be heatlhy")
	}
	if healthChecks["TestRegister"].Message != "This is just a test" {
		t.Errorf("Failed to register and get healthcheck status")
	}
}
func TestHealthCheckHandlerHealthy(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	gohealthy := Init()
	gohealthy.Register(mock{})
	gohealthy.HealthCheckHandler(response, request)
	expected := `This is just a test`
	if !strings.Contains(response.Body.String(), expected) {
		t.Errorf("Reponse body is %s, expected %s", response.Body, expected)
	}
	if response.Code != http.StatusOK {
		t.Errorf("Response code is %v, should be 200", response.Code)
	}
}

type mockUnhealthy struct{}

func (m mockUnhealthy) GetHealth() HealthStatus {
	return HealthStatus{"TestRegister", "This is just a test", false}
}
func TestHealthCheckHandlerUnhealthy(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	gohealthy := Init()
	gohealthy.Register(mockUnhealthy{})
	gohealthy.HealthCheckHandler(response, request)
	expected := `This is just a test`
	if !strings.Contains(response.Body.String(), expected) {
		t.Errorf("Reponse body is %s, expected %s", response.Body, expected)
	}
	if response.Code != http.StatusInternalServerError {
		t.Errorf("Response code is %v, should be 500", response.Code)
	}
}
