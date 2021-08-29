package main

import (
	"consistent-hashing/constant"
	"consistent-hashing/hashing"
	"flag"
	"github.com/golang/glog"
	"os"
)

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set(constant.LOG_TO_STD_ERR, constant.TRUE)
	flag.Set(constant.THRESHOLD, constant.WARN)
	flag.Set("v", "2")
	flag.Parse()
}

func performAllOperations(loadBalancer hashing.Hashing) {

	loadBalancer.AddNode("10.0.0.23")
	loadBalancer.AddNode("10.0.0.24")
	loadBalancer.AddNode("10.0.0.25")
	loadBalancer.AddNode("10.0.0.26")
	glog.Info(loadBalancer.GetTargetNode("213s1233-adff-asff-grgr-aldsbfbsqsdf"))
}

func main() {

	glog.Info("Starting our load balancer")

	loadBalancer := hashing.NewConsistentHashing()
	performAllOperations(loadBalancer)
}
