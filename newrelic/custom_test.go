package newrelic_test

import (
	"fmt"
	"github.com/inteleon/snap-plugin-collector-newrelic/newrelic"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	nr "github.com/yfronto/newrelic"
	"strings"
	"testing"
)

type customClientTestImpl struct {
	metricDataAppIDs         []int
	metricDataAppNames       map[int][]string
	metricDataComponentIDs   []int
	metricDataComponentNames map[int][]string
}

func (c *customClientTestImpl) GetApplicationMetricData(appID int, names []string, options *nr.MetricDataOptions) (*nr.MetricDataResponse, error) {
	c.metricDataAppIDs = append(c.metricDataAppIDs, appID)

	if len(c.metricDataAppNames) == 0 {
		c.metricDataAppNames = map[int][]string{}
	}

	c.metricDataAppNames[appID] = []string{}
	for i := range names {
		c.metricDataAppNames[appID] = append(c.metricDataAppNames[appID], names[i])
	}

	return &nr.MetricDataResponse{
		Metrics: []nr.MetricData{
			{
				Name: "hax",
				Timeslices: []nr.MetricTimeslice{
					{
						Values: map[string]float64{
							"average_response_time": 100.34,
							"throughput":            23,
						},
					},
				},
			},
		},
	}, nil
}

func (c *customClientTestImpl) GetComponentMetricData(componentID int, names []string, options *nr.MetricDataOptions) (*nr.MetricDataResponse, error) {
	c.metricDataComponentIDs = append(c.metricDataComponentIDs, componentID)

	if len(c.metricDataComponentNames) == 0 {
		c.metricDataComponentNames = map[int][]string{}
	}

	c.metricDataComponentNames[componentID] = []string{}
	for i := range names {
		c.metricDataComponentNames[componentID] = append(c.metricDataComponentNames[componentID], names[i])
	}

	return &nr.MetricDataResponse{
		Metrics: []nr.MetricData{
			{
				Name: "hacker",
				Timeslices: []nr.MetricTimeslice{
					{
						Values: map[string]float64{
							"average_response_time": 13.37,
							"throughput":            33,
						},
					},
				},
			},
		},
	}, nil
}

func TestGetCustomMetricTypesSuccess(t *testing.T) {
	c := &newrelic.Custom{}

	metrics, err := c.GetMetricTypes(plugin.Config{})

	if err != nil {
		t.Fatal(err)
	}

	expectedLen := len(newrelic.CustomMetrics)
	if len(metrics) != expectedLen {
		t.Fatal("expected", expectedLen, "got", len(metrics))
	}

	for i, m := range metrics {
		expectedNS := fmt.Sprintf("inteleon/newrelic/metric/%s", strings.Join(newrelic.CustomMetrics[i].Namespace.Strings(), "/"))
		ns := strings.Join(m.Namespace.Strings(), "/")
		t.Log(ns)

		if ns != expectedNS {
			t.Fatal("expected", expectedNS, "got", ns)
		}
	}
}

func TestCollectCustomMetricsSuccess(t *testing.T) {
	customClient := &customClientTestImpl{}

	c := &newrelic.Custom{
		CustomClient: customClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "application", "1337", "*", "hax", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "application",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "application", "1337", "1", "hax", "throughput", "value"),
			Tags: map[string]string{
				"Type": "application",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "component", "31337", "*", "hacker", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "component",
				"Unit": "float",
			},
		},
	}

	ret, err := c.CollectMetrics(metrics)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 3 {
		t.Fatal("expected", 3, "got", len(ret))
	}

	if ret[0].Data.(float64) != 100.34 {
		t.Fatal("expected", 100.34, "got", ret[0].Data.(float64))
	}

	if ret[1].Data.(float64) != 23 {
		t.Fatal("expected", 23, "got", ret[1].Data.(float64))
	}

	if ret[2].Data.(float64) != 13.37 {
		t.Fatal("expected", 13.37, "got", ret[2].Data.(float64))
	}

	if len(customClient.metricDataAppIDs) != 1 {
		t.Fatal("expected", 1, "got", len(customClient.metricDataAppIDs))
	}

	if customClient.metricDataAppIDs[0] != 1337 {
		t.Fatal("expected", 1337, "got", customClient.metricDataAppIDs[0])
	}

	if len(customClient.metricDataAppNames) != 1 {
		t.Fatal("expected", 1, "got", len(customClient.metricDataAppNames))
	}

	if len(customClient.metricDataAppNames[1337]) != 1 {
		t.Fatal("expected", 1, "got", len(customClient.metricDataAppNames[1337]))
	}

	if customClient.metricDataAppNames[1337][0] != "hax" {
		t.Fatal("expected", "hax", "got", customClient.metricDataAppNames[1337][0])
	}

	if len(customClient.metricDataComponentIDs) != 1 {
		t.Fatal("expected", 1, "got", len(customClient.metricDataComponentIDs))
	}

	if customClient.metricDataComponentIDs[0] != 31337 {
		t.Fatal("expected", 31337, "got", customClient.metricDataComponentIDs[0])
	}

	if len(customClient.metricDataComponentNames) != 1 {
		t.Fatal("expected", 1, "got", len(customClient.metricDataComponentNames))
	}

	if len(customClient.metricDataComponentNames[31337]) != 1 {
		t.Fatal("expected", 1, "got", len(customClient.metricDataComponentNames[31337]))
	}

	if customClient.metricDataComponentNames[31337][0] != "hacker" {
		t.Fatal("expected", "hacker", "got", customClient.metricDataComponentNames[31337][0])
	}
}

func TestCollectCustomMetricsMetricNameNotFoundFailure(t *testing.T) {
	customClient := &customClientTestImpl{}

	c := &newrelic.Custom{
		CustomClient: customClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "application", "1337", "*", "hax", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "application",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "application", "1337", "*", "h4x", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "application",
				"Unit": "float",
			},
		},
	}

	_, err := c.CollectMetrics(metrics)
	if err == nil {
		t.Fatal("expected", "error", "got", nil)
	}

	expectedErrStr := "Metric name mismatch! Requested metric name: h4x. Metric name in the received payload: hax."
	if err.Error() != expectedErrStr {
		t.Fatal("expected", expectedErrStr, "got", err.Error())
	}
}

func TestCollectCustomMetricsMetricValueNameNotFoundFailure(t *testing.T) {
	customClient := &customClientTestImpl{}

	c := &newrelic.Custom{
		CustomClient: customClient,
	}

	metrics := []plugin.Metric{
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "application", "1337", "*", "hax", "average_response_time", "value"),
			Tags: map[string]string{
				"Type": "application",
				"Unit": "float",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "metric", "application", "1337", "*", "hax", "h444x", "value"),
			Tags: map[string]string{
				"Type": "application",
				"Unit": "float",
			},
		},
	}

	_, err := c.CollectMetrics(metrics)
	if err == nil {
		t.Fatal("expected", "error", "got", nil)
	}

	expectedErrStr := "Path element not found: h444x"
	if err.Error() != expectedErrStr {
		t.Fatal("expected", expectedErrStr, "got", err.Error())
	}
}
