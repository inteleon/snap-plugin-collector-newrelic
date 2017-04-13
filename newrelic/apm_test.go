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
	appIDs []int
}

func (a *apmClientTestImpl) GetApplication(appID int) (*nr.Application, error) {
	a.appIDs = append(a.appIDs, appID)

	return &nr.Application{
		HealthStatus: "awesome",
		Reporting:    true,
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
		expectedNS := fmt.Sprintf("inteleon/newrelic/apm/*/%s", newrelic.APMMetrics[i]["namespace"])
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
				"namespace": "show/health/status",
				"type":      "application",
				"path":      "HealthStatus",
				"unit":      "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1234", "show", "health", "status"),
			Tags: map[string]string{
				"namespace": "show/health/status",
				"type":      "application",
				"path":      "HealthStatus",
				"unit":      "string",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "apm", "1337", "show", "reporting"),
			Tags: map[string]string{
				"namespace": "show/reporting",
				"type":      "application",
				"path":      "Reporting",
				"unit":      "bool",
			},
		},
		{
			Namespace: plugin.NewNamespace("inteleon", "newrelic", "browser", "314", "show", "health", "status"),
			Tags: map[string]string{
				"namespace": "show/health/status",
				"type":      "application",
				"path":      "HealthStatus",
				"unit":      "string",
			},
		},
	}

	ret, err := a.CollectMetrics(metrics)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 3 {
		t.Fatal("expected", 3, "got", len(ret))
	}

	for _, m := range ret[:2] {
		if m.Data.(string) != "awesome" {
			t.Fatal("expected", "awesome", "got", m.Data.(string))
		}
	}

	if !ret[2].Data.(bool) {
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
