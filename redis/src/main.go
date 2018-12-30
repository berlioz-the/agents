package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/berlioz-the/connector-go"
)

const redisConfFile = "/etc/redis.conf"

var configContents = make(map[string]string)
var isConfigReady = false
var isRedisUp = false
var isSelfDefaultPresent = false
var isSelfGossipPresent = false
var isClusterDeployed = false

var replicaCount = 0
var myIdentity = ""
var myService = ""

var setupMutex = &sync.Mutex{}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkConfigReady() {
	if isConfigReady {
		return
	}

	if isSelfDefaultPresent && isSelfGossipPresent {
		isConfigReady = true
	}

	if !isConfigReady {
		return
	}

	fmt.Printf("**** CONFIG IS READY: %#v\n", configContents)

	saveConfig()

	startRedis()
}

func saveConfig() {
	f, err := os.Create(redisConfFile)
	check(err)

	defer f.Close()

	for k, v := range configContents {
		f.WriteString(k + " " + v + "\n")
	}

	f.Sync()
}

func startRedis() {
	fmt.Printf("**** STARTING REDIS\n")

	var commandArray []string
	commandArray = append(commandArray, "docker-entrypoint.sh")
	commandArray = append(commandArray, "redis-server")
	commandArray = append(commandArray, redisConfFile)

	commandArray = append(commandArray, "--cluster-announce-ip")
	commandArray = append(commandArray, os.Getenv("BERLIOZ_ADDRESS"))
	commandArray = append(commandArray, "--cluster-announce-port")
	commandArray = append(commandArray, os.Getenv("BERLIOZ_PROVIDED_PORT_DEFAULT"))
	commandArray = append(commandArray, "--cluster-announce-bus-port")
	commandArray = append(commandArray, os.Getenv("BERLIOZ_PROVIDED_PORT_GOSSIP"))

	fmt.Printf("***** REDIS START COMMAND: %#v\n", commandArray)
	cmd := exec.Command(commandArray[0], commandArray[1:]...)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		fmt.Printf("**** ERROR STARTING REDIS: %#v\n", err)
		os.Exit(1)
	}
	fmt.Printf("**** STARTED REDIS. Pid: %d\n", cmd.Process.Pid)
	time.Sleep(5 * time.Second)
	fmt.Printf("**** REDIS START SLEEP COMPLETE> Pid: %d\n", cmd.Process.Pid)

	isRedisUp = true

	setupCluster(berlioz.Service(myService).All())
}

func setupCluster(peers map[string]interface{}) {
	if !isRedisUp {
		return
	}
	if isClusterDeployed {
		return
	}
	var requiredCount = (replicaCount + 1) * 3

	fmt.Printf("***** SETUP CLUSTER PEERS: %#v\n", peers)

	var orderedPeers []string
	for i := 1; i <= requiredCount; i++ {
		identity := strconv.Itoa(i)
		if peerObj, ok := peers[identity]; ok {
			peerDict := peerObj.(map[string]interface{})
			peerArg := fmt.Sprintf("%v:%v", peerDict["address"], peerDict["port"])
			orderedPeers = append(orderedPeers, peerArg)
		} else {
			return
		}
	}

	fmt.Printf("***** ORDERED PEERS: %#v\n", orderedPeers)

	go func() {
		setupMutex.Lock()
		defer setupMutex.Unlock()

		if isClusterDeployed {
			return
		}

		var commandArray []string
		commandArray = append(commandArray, "ruby")
		commandArray = append(commandArray, "/var/local/redis/redis-trib.rb")
		commandArray = append(commandArray, "create")
		commandArray = append(commandArray, "--verbose")
		commandArray = append(commandArray, "--replicas")
		commandArray = append(commandArray, fmt.Sprintf("%v", replicaCount))
		for _, peerArg := range orderedPeers {
			commandArray = append(commandArray, peerArg)
		}

		fmt.Printf("***** SETUP CLUSTER COMMAND: %#v\n", commandArray)

		fmt.Printf("**** CONFIGURING CLUSTER\n")
		cmd := exec.Command(commandArray[0], commandArray[1:]...)
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		if err != nil {
			fmt.Printf("**** ERROR CONFIGURING CLUSTER: %#v\n", err)
		}

		isClusterDeployed = true
	}()
}

func forever() {
	for {
		// fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}

func main() {

	myIdentity = berlioz.Identity()
	myService = os.Getenv("BERLIOZ_SERVICE")
	if val, err := strconv.Atoi(os.Getenv("replica_count")); err == nil {
		replicaCount = val
	}

	if myIdentity != "1" {
		isClusterDeployed = true
	}

	fmt.Printf("**** IDENTITY: %v\n", myIdentity)
	fmt.Printf("**** REPLICA COUNT: %v\n", replicaCount)

	configContents["port"] = "6379"
	configContents["cluster-enabled"] = "yes"
	configContents["cluster-config-file"] = "nodes.conf"
	configContents["cluster-node-timeout"] = "5000"
	configContents["appendonly"] = "yes"
	// configContents["pidfile"] = "/etc/redis.pid"
	// configContents["loglevel"] = "debug"
	// configContents["logfile"] = "/etc/redis-server.log"

	configContents["cluster-announce-ip"] = os.Getenv("BERLIOZ_ADDRESS")
	configContents["cluster-announce-port"] = os.Getenv("BERLIOZ_PROVIDED_PORT_DEFAULT")
	configContents["cluster-announce-bus-port"] = os.Getenv("BERLIOZ_PROVIDED_PORT_GOSSIP")
	isSelfDefaultPresent = true
	isSelfGossipPresent = true
	checkConfigReady()

	// berlioz.MyEndpoint("default").Monitor(func(ep berlioz.EndpointModel) {
	// 	fmt.Printf("**** MONITOR DEFAULT EP: %#v. Present: %t\n", ep, ep.IsPresent())
	// })

	// berlioz.MyEndpoint("gossip").Monitor(func(ep berlioz.EndpointModel) {
	// 	fmt.Printf("**** MONITOR GOSSIP EP: %#v. Present: %t\n", ep, ep.IsPresent())
	// })
	berlioz.Service(myService).MonitorAll(func(peers map[string]interface{}) {
		fmt.Printf("***** UPDATED DEFAULT PEERS: %#v\n", peers)
		setupCluster(peers)
	})

	forever()
}

// if !isSelfDefaultPresent {
// 	if selfPeerObj, ok := peers[myIdentity]; ok {
// 		fmt.Printf("***** UPDATED DEFAULT SELF PEER: %#v\n", selfPeerObj)

// 		selfPeer := selfPeerObj.(map[string]interface{})
// 		configContents["cluster-announce-ip"] = fmt.Sprintf("%v", selfPeer["address"])
// 		configContents["cluster-announce-port"] = fmt.Sprintf("%v", selfPeer["port"])

// 		isSelfDefaultPresent = true
// 		checkConfigReady()
// 	}
// }

// berlioz.Service(myService).Endpoint("gossip").MonitorAll(func(peers map[string]interface{}) {
// 	fmt.Printf("***** UPDATED GOSSIP PEERS: %#v\n", peers)
// 	// if !isSelfGossipPresent {
// 	// 	if selfPeerObj, ok := peers[myIdentity]; ok {
// 	// 		fmt.Printf("***** UPDATED GOSSIP SELF PEER: %#v\n", selfPeerObj)

// 	// 		selfPeer := selfPeerObj.(map[string]interface{})
// 	// 		configContents["cluster-announce-bus-port"] = fmt.Sprintf("%v", selfPeer["port"])

// 	// 		isSelfGossipPresent = true
// 	// 		checkConfigReady()
// 	// 	}
// 	// }

// 	setupCluster()
// })
