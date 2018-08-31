package main

import (
	"fmt"
	"os"
	"os/exec"
)

var runningCommands = []*exec.Cmd{}

func startGrafana() {
	fmt.Printf("**** STARTING GRAFANA\n")

	cmd := exec.Command("/run.sh")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		fmt.Printf("**** ERROR STARTING GRAFANA: %#v\n", err)
		os.Exit(1)
	}
	fmt.Printf("**** STARTED GRAFANA. Pid: %d\n", cmd.Process.Pid)
	runningCommands = append(runningCommands, cmd)
}

func stopGrafana() {
	for _, cmd := range runningCommands {
		fmt.Printf("**** KILLING GRAFANA. Pid: %d...\n", cmd.Process.Pid)
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("*** ERROR. Failed to kill process: %s, Error: %#v\n", cmd.Process.Pid, err)
		}
		fmt.Printf("**** GRAFANA KILLED.\n")
	}
	runningCommands = []*exec.Cmd{}
}

func restartGrafana() {
	stopGrafana()
	startGrafana()
}
