package gohealthy

import "testing"

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
