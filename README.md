# Go Sysdig Client Library 

A high level Go client library for Sysdig resources like:

 * GetSumMetric


## Installation

`go get -u github.com/NeowayLabs/sysdigcli`

## Quick Start

```go
package main


import "github.com/NeowayLabs/sysdigcli"


func main() {
	filter := "agent.tag.team = \"datapirates\" and bot.name = \"rf-pf-input\""
    
    	metrics := []sysdigcli.Metric{
    		{
    			Id: "bot.result.status.2xx.total",
    			Aggregations: sysdigcli.Aggregation{
    				Time: "sum",
    				Group: "sum",
    			},
    		},
    	}
    	period := sysdigcli.Period{
    		Days: 30,
    	}
    
    	s := sysdigcli.New("Bearer token_here")
    	
    	sumMetric, err := s.GetSumMetric(metrics, filter, period)
}
```

To see more examples, [check integration test file](sysdigcli_integration_test.go)

## Running tests

Integration tests (to execute the test you need put your token in the `setup` of `sysdigcli_integration_test.go`):

    make check-integration
    

## Analyze code

This target run this tools:

 * vet 
 * staticcheck 
 * gosimple
 * unused 

How to run:

    make analyze

## Release version

    make release version=(VERSION)