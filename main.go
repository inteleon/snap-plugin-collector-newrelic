package main

import (
	"github.com/inteleon/snap-plugin-collector-newrelic/newrelic"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	pluginName    = "newrelic-collector"
	pluginVersion = 1
)

func main() {
	plugin.StartCollector(
		&newrelic.Collector{},
		pluginName,
		pluginVersion,
	)
}
