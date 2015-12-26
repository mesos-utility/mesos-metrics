package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soarpenguin/mesos-metrics/cron"
	"github.com/soarpenguin/mesos-metrics/g"
	mhttp "github.com/soarpenguin/mesos-metrics/http"
)

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	handleVersion(*version)

	// global config
	g.ParseConfig(*cfg)
	g.InitRpcClients()

	cron.Collect()

	// http
	go mhttp.Start()

	select {}
}
