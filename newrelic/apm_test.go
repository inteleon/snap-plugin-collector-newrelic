package newrelic_test

import (
	"fmt"
	"github.com/inteleon/snap-plugin-collector-newrelic/newrelic"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	nr "github.com/yfronto/newrelic"
	"strings"
	"testing"
)

type apmClientTestImpl struct {
	appIDs           []int
	metricDataAppIDs []int
	metricDataNames  map[int][]string
}

func (a *apmClientTestImpl) GetApplication(appID int) (*nr.Application, error) {
	a.appIDs = append(a.appIDs, appID)

	return &nr.Application{
		ApplicationSummary: nr.ApplicationSummary{
			ResponseTime: 13.37,
		},
		HealthStatus: "awesome",
		Reporting:    true,
	}, nil
}

func TestGetAppMetricTypesSuccess(t *testing.T) {
	a := &newrelic.APM{}

	metrics, err := a.GetMetricTypes(plugin.Config{})

	if err != nil {
		t.Fatal(err)
	}

	expectedLen := len(newrelic.APMMetrics)
	if len(metrics) != expectedLen {
		t.Fatal("expected", expectedLen, "got", len(metrics))
	}

	for i, m := range metrics {
		expectedNS := fmt.Sprintf("inteleon/newrelic/apm/%s", strings.Join(newrelic.APMMetrics[i].Namespace.Strings(), "/"))
		ns := strings.Join(m.Namespace.Strings(), "/")
		t.Log(ns)

		if ns != expectedNS {
			t.Fatal("expected", expectedNS, "got", ns)
		}
	}
}

func TestCollectAppMetricsAppIDSuccess(t *testing.T) {
	apmClient := &apmClientTestImpl{}

	a := &newrelic.APM{
		APMClient: apmClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "application", "1337", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "application", "1234", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "application", "1337", "show", "summary", "application", "response_time"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "ApplicationSummary/ResponseTime",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "application", "1337", "show", "reporting"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "Reporting",
				"Unit": "bool",
			},
		},
	}

	ret, err := a.CollectMetrics(metrics)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 4 {
		t.Fatal("expected", 4, "got", len(ret))
	}

	for _, m := range ret[:2] {
		if m.Data.(string) != "awesome" {
			t.Fatal("expected", "awesome", "got", m.Data.(string))
		}
	}

	appResponseTimeData := ret[2].Data.(float64)
	if appResponseTimeData != 13.37 {
		t.Fatal("expected", 13.37, "got", appResponseTimeData)
	}

	if !ret[3].Data.(bool) {
		t.Fatal("expected", true, "got", false)
	}

	if len(apmClient.appIDs) != 2 {
		t.Fatal("expected", 2, "got", len(apmClient.appIDs))
	}

	expectedIDs := []int{1337, 1234}
	for i, id := range apmClient.appIDs {
		if id != expectedIDs[i] {
			t.Fatal("expected", expectedIDs[i], "got", id)
		}
	}
}
