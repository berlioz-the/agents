package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/berlioz-the/connector-go"
	yaml "gopkg.in/yaml.v2"
)

type serviceInfo struct {
	id        string
	serviceID string
	endpoint  string
	handler   func(map[string]interface{})
	monitor   berlioz.SubscribeInfo
}

var prometheusPid int = 0

var monitoredAgents = make(map[string]serviceInfo)
var newAgents = make(map[string]bool)
var trackedPeers = make(map[string]map[string]interface{})

func forever() {
	for {
		// fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}

func constructConfig() Config {
	config := Config{}
	config.GlobalConfig.ScrapeInterval = "10s"
	config.GlobalConfig.EvaluationInterval = "10s"

	scrapeStaticConfig := ScrapeStaticConfig{}
	scrapeStaticConfig.Targets = []string{}
	// scrapeConfig.StaticConfigs.

	scrapeConfig := ScrapeConfig{}
	scrapeConfig.JobName = "berlioz"
	scrapeConfig.StaticConfigs = append(scrapeConfig.StaticConfigs, &scrapeStaticConfig)
	config.ScrapeConfigs = append(config.ScrapeConfigs, &scrapeConfig)

	for _, peerMap := range trackedPeers {
		for _, rawPeer := range peerMap {
			fmt.Printf("***** [constructConfig] PEER: %#v\n", rawPeer)
			peer := rawPeer.(map[string]interface{})
			target := fmt.Sprintf("%s:%d", peer["address"], int(peer["port"].(float64)))
			fmt.Printf("***** [constructConfig] Target: %#v\n", target)

			scrapeStaticConfig.Targets = append(scrapeStaticConfig.Targets, target)
		}
	}

	return config
}

func writeConfig() {
	config := constructConfig()
	y, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	fmt.Println(string(y))

	f, err := os.Create("/etc/prometheus/prometheus.yml")
	if err != nil {
		panic(err)
	}

	f.WriteString(string(y))
	f.Sync()
}

func notifyPrometheus() {
	if prometheusPid > 0 {
		fmt.Printf("Sending SIGHUP to Prometheus\n")
		syscall.Kill(prometheusPid, syscall.SIGHUP)
	}
}

func processPeers() {
	fmt.Printf("***** TRACKED PEERS: %#v\n", trackedPeers)
	writeConfig()
	notifyPrometheus()
}

func monitorAgent(consumed berlioz.ConsumesModel) {
	fmt.Printf("***** AGENT TO MONITOR: %s\n", consumed.ID)
	id := consumed.ID + "-" + consumed.Endpoint
	newAgents[id] = true

	if _, ok := monitoredAgents[id]; ok {
		return
	}

	info := serviceInfo{id: id, serviceID: consumed.ID, endpoint: consumed.Endpoint}
	info.handler = func(peers map[string]interface{}) {
		fmt.Printf("***** PEERS CHANGED FOR: %s\n", id)
		trackedPeers[id] = peers
		processPeers()
	}

	monitoredAgents[id] = info

	info.monitor = berlioz.Sector(consumed.Sector).Service(consumed.Name).Endpoint(consumed.Endpoint).MonitorAll(info.handler)
}

func stopMonitoring(id string, serviceInfo serviceInfo) {
	fmt.Printf("***** AGENT TO STOP MONITORING: %s\n", id)
	serviceInfo.monitor.Stop()
	delete(trackedPeers, id)
	delete(monitoredAgents, id)
	processPeers()
}

func applyAgentChanges() {
	// fmt.Printf("***** applyAgentChanges: %#v\n", monitoredAgents)

	for id, serviceInfo := range monitoredAgents {
		if _, ok := newAgents[id]; !ok {
			stopMonitoring(id, serviceInfo)
		}
	}
}

func onConsumesChanged(consumes []berlioz.ConsumesModel) {
	// fmt.Printf("***** UPDATED MONITOR CONSUMES: %#v\n", consumes)
	newAgents = make(map[string]bool)
	for _, consumed := range consumes {
		if consumed.Kind == "service" && consumed.Name == "berlioz_agent" && consumed.Endpoint == "mon" {
			monitorAgent(consumed)
		}
	}
	applyAgentChanges()
}

func main() {
	var argsWithoutProg = os.Args[1:]
	if len(argsWithoutProg) > 0 {
		pid, err := strconv.Atoi(argsWithoutProg[0])
		if err != nil {
			fmt.Printf("ERROR. Invalid prometheus pid provided: %#v\n", argsWithoutProg)
			os.Exit(1)
		}
		prometheusPid = pid
	}

	fmt.Printf("prometheusPid: %s\n", prometheusPid)

	berlioz.Consumes().MonitorAll(onConsumesChanged)

	forever()
}
