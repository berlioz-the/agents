package main

import (
	"net/url"
)

type LabelSet map[string]string
type Duration string

// Config is the top-level configuration for Prometheus's config files.
type Config struct {
	GlobalConfig GlobalConfig `yaml:"global"`
	// AlertingConfig AlertingConfig  `yaml:"alerting,omitempty"`
	// RuleFiles      []string        `yaml:"rule_files,omitempty"`
	ScrapeConfigs []*ScrapeConfig `yaml:"scrape_configs,omitempty"`

	// RemoteWriteConfigs []*RemoteWriteConfig `yaml:"remote_write,omitempty"`
	// RemoteReadConfigs  []*RemoteReadConfig  `yaml:"remote_read,omitempty"`
}

// GlobalConfig configures values that are used across other configuration
// objects.
type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval Duration `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout Duration `yaml:"scrape_timeout,omitempty"`
	// How frequently to evaluate rules by default.
	EvaluationInterval Duration `yaml:"evaluation_interval,omitempty"`
	// The labels to add to any timeseries that this Prometheus instance scrapes.
	ExternalLabels LabelSet `yaml:"external_labels,omitempty"`
}

// ScrapeConfig configures a scraping unit for Prometheus.
type ScrapeConfig struct {
	// The job name to which the job label is set by default.
	JobName string `yaml:"job_name"`
	// Indicator whether the scraped metrics should remain unmodified.
	HonorLabels bool `yaml:"honor_labels,omitempty"`
	// A set of query parameters with which the target is scraped.
	Params url.Values `yaml:"params,omitempty"`
	// How frequently to scrape the targets of this scrape config.
	ScrapeInterval Duration `yaml:"scrape_interval,omitempty"`
	// The timeout for scraping targets of this config.
	ScrapeTimeout Duration `yaml:"scrape_timeout,omitempty"`
	// The HTTP resource path on which to fetch metrics from targets.
	MetricsPath string `yaml:"metrics_path,omitempty"`
	// The URL scheme with which to fetch metrics from targets.
	Scheme string `yaml:"scheme,omitempty"`
	// More than this many samples post metric-relabelling will cause the scrape to fail.
	SampleLimit uint `yaml:"sample_limit,omitempty"`

	StaticConfigs []*ScrapeStaticConfig `yaml:"static_configs,omitempty"`

	// We cannot do proper Go type embedding below as the parser will then parse
	// values arbitrarily into the overflow maps of further-down types.

	// ServiceDiscoveryConfig sd_config.ServiceDiscoveryConfig `yaml:",inline"`
	// HTTPClientConfig       config_util.HTTPClientConfig     `yaml:",inline"`

	// List of target relabel configurations.
	// RelabelConfigs []*RelabelConfig `yaml:"relabel_configs,omitempty"`
	// List of metric relabel configurations.
	// MetricRelabelConfigs []*RelabelConfig `yaml:"metric_relabel_configs,omitempty"`
}

// ScrapeStaticConfig
type ScrapeStaticConfig struct {
	Targets []string `yaml:"targets,omitempty"`
}
