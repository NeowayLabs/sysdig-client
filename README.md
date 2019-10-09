# Go Sysdig Client Library 


[![Build Status](https://travis-ci.org/NeowayLabs/gootstrap.svg?branch=master)](https://travis-ci.org/NeowayLabs/sysdig-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/NeowayLabs/sysdig-client)](https://goreportcard.com/report/github.com/NeowayLabs/sysdig-client)

A high level Go client library for Sysdig resources like:

 * GetSumMetric
 

More details about the official sysdig API can be found [here](https://sysdig.gitbooks.io/sysdig-cloud-api/content/rest_api/data.html)


## Prerequisite

To use this lib you need set your sysdig API token in the var env `SYSDIG_CLOUD_API_TOKEN`. You can find your token [here](https://app.sysdigcloud.com/#/settings/user)

## Installation

`go get -u github.com/NeowayLabs/sysdig-client`

## Quick Start

```go
package main


import sysdigclient "github.com/NeowayLabs/sysdig-client"


func main() {
	filter := "agent.tag.team = \"datapirates\" and bot.name = \"abi\""
    
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
    
    	s := sysdigclient.New()
    	
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