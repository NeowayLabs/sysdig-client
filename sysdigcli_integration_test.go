// +build integration

package sysdigcli_test

import (
	"github.com/NeowayLabs/sysdigcli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSumMetric(t *testing.T) {

	s := setup(t)

	filter := "agent.tag.team = \"datapirates\" and bot.name = \"rf-pf-input\""

	metrics := []sysdigcli.Metric{
		{
			Id: "bot.result.status.2xx.total",
			Aggregations: sysdigcli.Aggregation{
				Time:  "sum",
				Group: "sum",
			},
		},
	}
	period := sysdigcli.Period{
		Days: 30,
	}

	sumMetric, err := s.GetSumMetric(metrics, filter, period)

	assert.Nil(t, err)

	assert.Greater(t, sumMetric, 0, "Count should be greater than 0!")
}

func setup(t *testing.T) *sysdigcli.Sysdigcli {
	return sysdigcli.New("Bearer token_here")
}
