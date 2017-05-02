# Snap Telemetry New Relic Collector

[![GoDoc](https://godoc.org/github.com/inteleon/snap-plugin-collector-newrelic/newrelic?status.svg)](https://godoc.org/github.com/inteleon/snap-plugin-collector-newrelic/newrelic)

## What the plugin does

This plugin collect metrics from `New Relic`.

It's a very early release so there's a lot of missing functionality.

See [Available metrics](#available-metrics) for more information about what is supported.

## Supported platforms

Should work on any platform that `Snap` supports.

## Known issues

### grpc

The latest `grpc` dependency does not work with the latest `plugin toolkit` for `snaptel`.

See `tools/ci/setup.sh` for a "solution".

### New Relic Go library

I created a pull request with improvements to the `New Relic` Go library that I'm using.
So far these changes are still a pull request, not merged into the master branch.

Therefore you must clone my `New Relic` Go library repository and check out the correct branch.

```bash
mkdir -p $GOPATH/src/github.com/yfronto
git clone https://github.com/inteleon/newrelic.git $GOPATH/src/github.com/yfronto/newrelic
cd $GOPATH/src/github.com/yfronto/newrelic
git checkout feature/support-for-component-metrics
```

## Snap version dependencies

Developed and tested with `Snap` version `1.2.0`.

## Installation

Download and include the binaries in your `Snap` plugin directory.

## Usage

### Available metrics

Currently we only support application APM and component metrics

You can fetch all basic metrics for your application and also more specified metrics, like external services, metrics from plugins, etc.

You can fetch a list of available metrics per application at https://rpm.newrelic.com/api/explore/applications/metric_names. Use these metric names to create a metric collection namespace in your configuration file.

Example:

```yaml
---
version: 1
schedule:
  type: "simple"
  interval: "60s"
deadline: "15s"
workflow:
  collect:
    metrics:
      /inteleon/newrelic/apm/application/APP_ID/show/summary/application/response_time: {}
      /inteleon/newrelic/apm/application/APP_ID/show/summary/application/throughput: {}
      /inteleon/newrelic/apm/application/APP_ID/show/summary/application/error_rate: {}
      "|inteleon|newrelic|metric|application|APP_ID|*|External/api.github.com/all|average_response_time|value": {} # Average value for the last 30 minutes (default New Relic timeframe).
      "|inteleon|newrelic|metric|application|APP_ID|1|External/api.github.com/all|calls_per_minute|value": {} # Average value for the last minute.
      "|inteleon|newrelic|metric|application|APP_ID|5|External/api.github.com/all|standard_deviation|value": {} # Average value for the last 5 minutes.
      "|inteleon|newrelic|metric|component|APP_ID|1|Component/Runtime/System/Threads[Threads]|average_value|value": {} # Number of threads a Go service is using (fetches using the GoRelic New Relic plugin).
    config:
      /inteleon/newrelic:
        api_key: "SUPER SECRET NEW RELIC API KEY"
    publish:
      -
        plugin_name: "cloudwatch"
        config:
          region: "eu-central-1"
          namespace: "NewRelic"
```

It's important to use `|` as a delimiter when fetching metrics, since most, or all, use `/` as part of the metric name.

### Example configuration

See [newrelic.example.yml](newrelic.example.yml) for a configuration example with all available metrics and configuration options.

**NOTE:** Not all additional application metrics are included in the example configuration file, because those can be generated dynamically depending on the available metric values.

### Good to know

The plugin always fetches the default time frame, which is the last 30 minutes. I plan on supporting relative timeframes.

## Contributors

Coming soon.

## License

Apache-2.0 - https://github.com/inteleon/snap-plugin-collector-newrelic/blob/master/LICENSE
