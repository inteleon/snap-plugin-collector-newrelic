package newrelic

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	nr "github.com/yfronto/newrelic"
	"strconv"
)

// APMMetrics is a list containing the available APM metrics and their properties.
var APMMetrics = []Metric{
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("health"),
			plugin.NewNamespaceElement("status"),
		},
		Type: "application",
		Path: "HealthStatus",
		Unit: "string",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("reporting"),
		},
		Type: "application",
		Path: "Reporting",
		Unit: "bool",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("response_time"),
		},
		Type: "application",
		Path: "ApplicationSummary/ResponseTime",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("throughput"),
		},
		Type: "application",
		Path: "ApplicationSummary/Throughput",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("error_rate"),
		},
		Type: "application",
		Path: "ApplicationSummary/ErrorRate",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("apdex_target"),
		},
		Type: "application",
		Path: "ApplicationSummary/ApdexTarget",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("apdex_score"),
		},
		Type: "application",
		Path: "ApplicationSummary/ApdexScore",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("host_count"),
		},
		Type: "application",
		Path: "ApplicationSummary/HostCount",
		Unit: "int",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("application"),
			plugin.NewNamespaceElement("instance_count"),
		},
		Type: "application",
		Path: "ApplicationSummary/InstanceCount",
		Unit: "int",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("user"),
			plugin.NewNamespaceElement("response_time"),
		},
		Type: "application",
		Path: "EndUserSummary/ResponseTime",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("user"),
			plugin.NewNamespaceElement("throughput"),
		},
		Type: "application",
		Path: "EndUserSummary/Throughput",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("user"),
			plugin.NewNamespaceElement("apdex_target"),
		},
		Type: "application",
		Path: "EndUserSummary/ApdexTarget",
		Unit: "float",
	},
	{
		Namespace: plugin.Namespace{
			plugin.NewNamespaceElement("application"),
			plugin.NamespaceElement{
				Name:        "app_id",
				Description: "Application id",
				Value:       "*",
			},
			plugin.NewNamespaceElement("show"),
			plugin.NewNamespaceElement("summary"),
			plugin.NewNamespaceElement("user"),
			plugin.NewNamespaceElement("apdex_score"),
		},
		Type: "application",
		Path: "EndUserSummary/ApdexScore",
		Unit: "float",
	},
}

// APMClient is the interface every AMP client needs to implement.
type APMClient interface {
	GetApplication(int) (*nr.Application, error)
}

// APMClientImpl is a real implementation of an APMClient.
type APMClientImpl struct {
	APIKey string
}

// GetApplication fetches application information from New Relic (APM).
func (a *APMClientImpl) GetApplication(appID int) (*nr.Application, error) {
	c := nr.NewClient(a.APIKey)

	return c.GetApplication(appID)
}

// APM represents the APM service part of New Relic.
type APM struct {
	APMClient APMClient
}

// NewAPM creates and returns a new APM object with a configured APMClient.
func NewAPM(apiKey string) Service {
	return &APM{
		APMClient: &APMClientImpl{
			APIKey: apiKey,
		},
	}
}

// GetMetricTypes returns the available APM metric types.
func (a *APM) GetMetricTypes(_ plugin.Config) ([]plugin.Metric, error) {
	ns := plugin.NewNamespace("inteleon", "newrelic", "apm")

	return metricTypes(ns, APMMetrics)
}

// CollectMetrics fetches the requested APM metrics and returns them.
func (a *APM) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	collectedMetrics := []plugin.Metric{}

	if len(metrics) == 0 {
		return collectedMetrics, fmt.Errorf("List of metrics is empty")
	}

	apps := []plugin.Metric{}
	for i, m := range metrics {
		if m.Namespace.Element(2).Value != "apm" {
			continue
		}

		apps = append(apps, metrics[i])
	}

	if len(apps) > 0 {
		appsMetrics, err := a.collectApplications(apps)
		if err != nil {
			return collectedMetrics, err
		}

		for i := range appsMetrics {
			collectedMetrics = append(collectedMetrics, appsMetrics[i])
		}
	}

	return collectedMetrics, nil
}

func (a *APM) collectApplications(metrics []plugin.Metric) ([]plugin.Metric, error) {
	appsMetrics := []plugin.Metric{}

	apps := map[int]*nr.Application{}
	for i, m := range metrics {
		appID := m.Namespace.Element(4)

		appIDInt, err := strconv.Atoi(appID.Value)
		if err != nil {
			return appsMetrics, err
		}

		if _, ok := apps[appIDInt]; !ok {
			// Application info missing, fetching...
			app, err := a.APMClient.GetApplication(appIDInt)
			if err != nil {
				return appsMetrics, err
			}

			apps[appIDInt] = app
		}

		// Convert the app data to a struct so it's more easily traversable and more universal before passing it to the populateMetric function.
		appMetric, err := populateMetric(metrics[i], structs.Map(apps[appIDInt]))
		if err != nil {
			// Metric not found, skip reporting it and continue execution.
			continue
		}

		appsMetrics = append(appsMetrics, appMetric)
	}

	return appsMetrics, nil
}
