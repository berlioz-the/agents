package main

import (
	"fmt"
	"time"

	"github.com/berlioz-the/connector-go"
)

func forever() {
	for {
		// fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}

func onPrometheusPeersChanged(peers map[string]interface{}) {
	config := constructDataSourceConfig(peers)
	writeYamlConfig(config, "/etc/grafana/provisioning/datasources/berlioz.yml")
	restartGrafana()
}

func onConsumesChanged(consumes []berlioz.ConsumesModel) {
	dashboard := constructDashboard(consumes)
	writeJsonConfig(dashboard, "/var/lib/grafana/dashboards/berlioz.json")
}

func main() {
	fmt.Printf("**** LAUNCHER STARTED\n")

	berlioz.Service("prometheus").Endpoint("server").MonitorAll(onPrometheusPeersChanged)
	berlioz.Consumes().MonitorAll(onConsumesChanged)

	startGrafana()

	forever()
}
