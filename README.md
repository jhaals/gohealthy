# gohealthy
[![Build Status](https://travis-ci.org/jhaals/gohealthy.svg?branch=master)](https://travis-ci.org/jhaals/gohealthy)

goheatlhy is a library for registering HTTP health checks that can be verified by a load balancer or service discovery. A HTTP 500 error will be returned if one or more of the health checks are failing.


## Example

    type customerPortal struct{}

    func (f customerPortal) GetHealth() gohealthy.HealthStatus {
     return gohealthy.HealthStatus{Name: "customer-satisfaction",
       Message: "Customer Satisfaction is above 180%",
       Healthy: true}
    }

    func main() {
     g := gohealthy.Init()
     g.Register(customerPortal{})
     g.RunServer(1337)
    }
