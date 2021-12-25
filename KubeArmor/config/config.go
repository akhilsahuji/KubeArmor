// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package config

import (
	"os"

	"flag"

	kg "github.com/kubearmor/KubeArmor/KubeArmor/log"
	"github.com/spf13/viper"
)

type KubearmorConfig struct {
	Grpc           string // GRPC Port to use
	Cluster        string // Cluster name to use for feeds
	LogPath        string // Log file to use
	HostVisibility string // Host visibility to use for kubearmor in process mode
	Policy         bool   // Enable/Disable policy enforcement
	HostPolicy     bool   // Enable/Disable host policy enforcement
	KVMAgent       bool   // Enable/Disable KVM Agent
}

var GlobalCfg KubearmorConfig

const CFG_CLUSTER string = "cluster"
const CFG_GRPC string = "gRPC"
const CFG_LOGPATH string = "logPath"
const CFG_HOSTVISIBILITY string = "hostVisibility"
const CFG_KUBEARMORPOLICY string = "enableKubeArmorPolicy"
const CFG_KUBEARMORHOSTPOLICY string = "enableKubeArmorHostPolicy"
const CFG_KUBEARMORVM string = "enableKubeArmorVm"

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func readCmdLineParams() {
	// Read configuration from command line
	clusterStr := flag.String(CFG_CLUSTER, "", "cluster name")
	grpcStr := flag.String(CFG_GRPC, "", "gRPC port number")
	logStr := flag.String(CFG_LOGPATH, "", "log file path, {path|stdout|none}")
	policyB := flag.Bool(CFG_KUBEARMORPOLICY, true, "enabling KubeArmorPolicy")
	hostPolicyB := flag.Bool(CFG_KUBEARMORHOSTPOLICY, false, "enabling KubeArmorHostPolicy")
	kvmAgentB := flag.Bool(CFG_KUBEARMORVM, false, "enabling KubeArmorVM")
	hostVisStr := flag.String(CFG_HOSTVISIBILITY, "process,file,network,capabilities", "Host Visibility to use [process,file,network,capabilities,none]")

	flag.Parse()
	if isFlagPassed(CFG_CLUSTER) {
		viper.Set(CFG_CLUSTER, *clusterStr)
	}
	if isFlagPassed(CFG_GRPC) {
		viper.Set(CFG_GRPC, *grpcStr)
	}
	if isFlagPassed(CFG_LOGPATH) {
		viper.Set(CFG_LOGPATH, *logStr)
	}
	if isFlagPassed(CFG_HOSTVISIBILITY) {
		viper.Set(CFG_HOSTVISIBILITY, *hostVisStr)
	}
	if isFlagPassed(CFG_KUBEARMORPOLICY) {
		viper.Set(CFG_KUBEARMORPOLICY, *policyB)
	}
	if isFlagPassed(CFG_KUBEARMORHOSTPOLICY) {
		viper.Set(CFG_KUBEARMORHOSTPOLICY, *hostPolicyB)
	}
	if isFlagPassed(CFG_KUBEARMORVM) {
		viper.Set(CFG_KUBEARMORVM, *kvmAgentB)
	}
}

func LoadConfig() error {
	// Set defaults
	viper.SetDefault(CFG_CLUSTER, "default")
	viper.SetDefault(CFG_GRPC, "32767")
	viper.SetDefault(CFG_LOGPATH, "/tmp/kubearmor.log")
	viper.SetDefault(CFG_KUBEARMORPOLICY, true)
	viper.SetDefault(CFG_KUBEARMORHOSTPOLICY, false)
	viper.SetDefault(CFG_KUBEARMORVM, false)
	viper.SetDefault(CFG_HOSTVISIBILITY, "process,file,network,capabilities")

	// Read configuration from env var
	// Note that the env var has to be set in uppercase for e.g, CLUSTER=xyz ./kubearmor
	viper.AutomaticEnv()

	// Read configuration from config file
	cfgfile := os.Getenv("KUBEARMOR_CFG")
	if cfgfile == "" {
		cfgfile = "kubearmor.yaml"
	}
	if _, err := os.Stat(cfgfile); err == nil {
		kg.Printf("setting config from file [%s]", cfgfile)
		viper.SetConfigFile(cfgfile)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}

	readCmdLineParams()

	GlobalCfg.Grpc = viper.GetString(CFG_GRPC)
	GlobalCfg.Cluster = viper.GetString(CFG_CLUSTER)
	GlobalCfg.LogPath = viper.GetString(CFG_LOGPATH)
	GlobalCfg.Policy = viper.GetBool(CFG_KUBEARMORPOLICY)
	GlobalCfg.HostPolicy = viper.GetBool(CFG_KUBEARMORHOSTPOLICY)
	GlobalCfg.KVMAgent = viper.GetBool(CFG_KUBEARMORVM)
	GlobalCfg.HostVisibility = viper.GetString(CFG_HOSTVISIBILITY)

	if GlobalCfg.KVMAgent {
		GlobalCfg.HostPolicy = true
	}

	kg.Printf("config [%+v]", GlobalCfg)
	return nil
}
