package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             int    `json:"port"`
	ReplicaOfHost    string `json:"replicaofhost"`
	ReplicaOfPort    int    `json:"replicaofport"`
	Role             string `json:"role"`
	ConnectedSlaves  int    `json:"connected_slaves"`
	MasterReplID     string `json:"master_replid"`
	MasterReplOffset int    `json:"master_repl_offset"`
}

func New() *Config {
	// Read command line arguments
	args := os.Args

	role := "master"
	port := 6379
	var replicaofHost string
	var replicaofPort int

	for i := 1; i < len(args); i++ {
		if args[i] == "--port" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &port)
		} else if args[i] == "--replicaof" && i+2 < len(args) {
			replicaofHost = args[i+1]
			fmt.Sscanf(args[i+2], "%d", &replicaofPort)
			role = "slave"
		}
	}

	config := &Config{
		Port:             port,
		Role:             role,
		ReplicaOfHost:    replicaofHost,
		ReplicaOfPort:    replicaofPort,
		ConnectedSlaves:  0,
		MasterReplID:     "123",
		MasterReplOffset: 0,
	}

	fmt.Println("Config created with: ", config)

	return config
}
