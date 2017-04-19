# Snap Telemetry New Relic Collector

[![GoDoc](https://godoc.org/github.com/inteleon/snap-plugin-collector-newrelic/newrelic?status.svg)](https://godoc.org/github.com/inteleon/snap-plugin-collector-newrelic/newrelic)

## What the plugin does

This plugin collect metrics from `New Relic`.

It's a very early release so there's a lot of missing functionality.

See [Available metrics](#available-metrics) for more information about what is supported.

## Supported platforms

Should work on any platform that `Snap` supports.

## Known issues

The latest `grpc` dependency does not work with the latest `snaptel plugin toolkit`.

## Snap version dependencies

Developed and tested with `Snap` version `1.2.0`.

## Installation

Coming soon.

## Usage

### Available metrics

Currently we only support application APM metrics.

You can fetch all basic metrics for your application and also more specified metrics, like external services, etc.

You can fetch a list of available metrics per application at https://rpm.newrelic.com/api/explore/applications/metric_names.

It's important to use `|` as a delimiter when fetching metrics, since most, or all, use `/` as part of the metric name.

### Example configuration

See [newrelic.example.yml](newrelic.example.yml) for a configuration example with all available metrics and configuration options.

## Contributors

Coming soon.

## License

Apache-2.0 - https://github.com/inteleon/snap-plugin-collector-newrelic/blob/master/LICENSE
