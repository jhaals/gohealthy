package gohealthy

import "testing"

type mock struct{}

func (m mock) Status() HealthStatus {
	return HealthStatus{"TestRegister", "This is just a test", true}
}
func TestRegister(t *testing.T) {

	goheatlhy := HealthChecks{}
	goheatlhy.Register(mock{})

	result, status := goheatlhy.Status()
	if !status {
		t.Errorf("this should be heatlhy")
	}
	if result["TestRegister"].Message != "This is just a test" {
		t.Errorf("Failed to register and get healthcheck status")
	}
}
