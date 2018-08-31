package main

import ()

// datasources:
// -  access: 'proxy'                       # make grafana perform the requests
//    editable: true                        # whether it should be editable
//    is_default: true                      # whether this should be the default DS
//    name: 'prom1'                         # name of the datasource
//    org_id: 1                             # id of the organization to tie this datasource to
//    type: 'prometheus'                    # type of the data source
//    url: 'http://172.17.0.6:9090'         # url of the prom instance
//    version: 1                            # well, versioning

// TBD
type ConfigDataSources struct {
	DataSources []*ConfigDataSource `yaml:"datasources,omitempty"`
}

// TBD
type ConfigDataSource struct {
	Access    string `yaml:"access,omitempty"`
	Editable  bool   `yaml:"editable,omitempty"`
	IsDefault bool   `yaml:"is_default,omitempty"`
	Name      string `yaml:"name,omitempty"`
	OrigID    int    `yaml:"org_id,omitempty"`
	Type      string `yaml:"type,omitempty"`
	URL       string `yaml:"url,omitempty"`
	Version   int    `yaml:"version,omitempty"`
}
