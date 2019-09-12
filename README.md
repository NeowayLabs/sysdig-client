# Go Sysdig Client Library 


[![Build Status](https://travis-ci.org/NeowayLabs/gootstrap.svg?branch=master)](https://travis-ci.org/NeowayLabs/sysdig-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/NeowayLabs/sysdig-client)](https://goreportcard.com/report/github.com/NeowayLabs/sysdig-client)

A high level Go client library for Sysdig resources like:

 * GetSumMetric


## Installation

`go get -u github.com/NeowayLabs/sysdig-client`

## Quick Start

```go
package main


import sysdigclient "github.com/NeowayLabs/sysdig-client"


func main() {
	filter := "agent.tag.team = \"datapirates\" and bot.name = \"rf-pf-input\""
    
    	metrics := []sysdigclient.Metric{
    		{
    			Id: "bot.result.status.2xx.total",
    			Aggregations: sysdigclient.Aggregation{
    				Time: "sum",
    				Group: "sum",
    			},
    		},
    	}
    	period := sysdigclient.Period{
    		Days: 30,
    	}
    
    	s := sysdigclient.New("Bearer token_here")
    	
    	sumMetric, err := s.GetSumMetric(metrics, filter, period)
}
```

To see more examples, [check integration test file](sysdigcli_integration_test.go)


## Running tests

Example:

    make check
    

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