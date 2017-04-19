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
		HealthStatus: "awesome",
		Reporting:    true,
	}, nil
}

func (a *apmClientTestImpl) GetApplicationMetricData(appID int, names []string, options *nr.MetricDataOptions) (*nr.MetricDataResponse, error) {
	a.metricDataAppIDs = append(a.metricDataAppIDs, appID)

	if len(a.metricDataNames) == 0 {
		a.metricDataNames = map[int][]string{}
	}

	a.metricDataNames[appID] = []string{}
	for i := range names {
		a.metricDataNames[appID] = append(a.metricDataNames[appID], names[i])
	}

	return &nr.MetricDataResponse{
		Metrics: []nr.MetricData{
			{
				Name: "hax",
				Timeslices: []nr.MetricTimeslice{
					{
						Values: map[string]float64{
							"average_response_time": 100.34,
						},
					},
				},
			},
		},
	}, nil
}

func TestGetMetricTypesSuccess(t *testing.T) {
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
		expectedNS := fmt.Sprintf("inteleon/newrelic/apm/*/%s", strings.Join(newrelic.APMMetrics[i].Namespace.Strings(), "/"))
		ns := strings.Join(m.Namespace.Strings(), "/")

		if ns != expectedNS {
			t.Fatal("expected", expectedNS, "got", ns)
		}
	}
}

func TestCollectMetricsAppIDSuccess(t *testing.T) {
	apmClient := &apmClientTestImpl{}

	a := &newrelic.APM{
		APMClient: apmClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1234", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "reporting"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "Reporting",
				"Unit": "bool",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "metric", "hax", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "metric",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "browser", "314", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
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

	if !ret[2].Data.(bool) {
		t.Fatal("expected", true, "got", false)
	}

	if ret[3].Data.(float64) != 100.34 {
		t.Fatal("expected", 100.34, "got", ret[3].Data.(float64))
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

	if len(apmClient.metricDataAppIDs) != 1 {
		t.Fatal("expected", 1, "got", len(apmClient.metricDataAppIDs))
	}

	if apmClient.metricDataAppIDs[0] != 1337 {
		t.Fatal("expected", 1337, "got", apmClient.metricDataAppIDs[0])
	}

	if len(apmClient.metricDataNames) != 1 {
		t.Fatal("expected", 1, "got", len(apmClient.metricDataNames))
	}

	if len(apmClient.metricDataNames[1337]) != 1 {
		t.Fatal("expected", 1, "got", len(apmClient.metricDataNames[1337]))
	}

	if apmClient.metricDataNames[1337][0] != "hax" {
		t.Fatal("expected", "hax", "got", apmClient.metricDataNames[1337][0])
	}
}

func TestCollectMetricsAppIDMetricNameNotFoundFailure(t *testing.T) {
	apmClient := &apmClientTestImpl{}

	a := &newrelic.APM{
		APMClient: apmClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1234", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "reporting"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "Reporting",
				"Unit": "bool",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "metric", "hax", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "metric",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "metric", "h4x", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "metric",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "browser", "314", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
	}

	_, err := a.CollectMetrics(metrics)
	if err == nil {
		t.Fatal("expected", "error", "got", nil)
	}

	expectedErrStr := "Metric name mismatch! Requested metric name: h4x. Metric name in the received payload: hax."
	if err.Error() != expectedErrStr {
		t.Fatal("expected", expectedErrStr, "got", err.Error())
	}
}

func TestCollectMetricsAppIDMetricValueNameNotFoundFailure(t *testing.T) {
	apmClient := &apmClientTestImpl{}

	a := &newrelic.APM{
		APMClient: apmClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1234", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "reporting"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "Reporting",
				"Unit": "bool",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "metric", "hax", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "metric",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "metric", "hax", "h444x", "value"),
			Tags: map[string]string{
				"Type": "metric",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "browser", "314", "show", "health", "status"),
			Tags: map[string]string{
				"Type": "application",
				"Path": "HealthStatus",
				"Unit": "string",
			},
		},
	}

	_, err := a.CollectMetrics(metrics)
	if err == nil {
		t.Fatal("expected", "error", "got", nil)
	}

	expectedErrStr := "Path element not found: h444x"
	if err.Error() != expectedErrStr {
		t.Fatal("expected", expectedErrStr, "got", err.Error())
	}
}
