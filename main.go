package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/open-falcon/common/model"
	"github.com/soarpenguin/mesos-metrics/g"
	//mhttp "github.com/soarpenguin/mesos-metrics/http"
)

func Hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := false

	if debug {
		log.Printf("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	var resp model.TransferResponse
	//err := TransferClient.Call("Transfer.Update", metrics, &resp)
	//if err != nil {
	//	log.Println("call Transfer.Update fail", err)
	//}

	if debug {
		log.Println("<=", &resp)
	}
}

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

	// http
	//mhttp.Start()

	//select {}
	gcfg := g.Config()
	addr := gcfg.Master.Apiurl

	var interval int64 = 6
	for {
	REST:
		time.Sleep(time.Duration(interval) * time.Second)
		resp, err := http.Get(addr)
		if err != nil {
			goto REST
			//panic(err)
		}
		defer resp.Body.Close()

		// read json http response
		jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			goto REST
			//panic(err)
		}

		var f interface{}
		err = json.Unmarshal(jsonDataFromHttp, &f)
		if err != nil {
			goto REST
			//panic(err)
		}

		now := time.Now().Unix()
		hostname, _ := Hostname()
		m := f.(map[string]interface{})
		for k, v := range m {
			key := fmt.Sprintf("mesos.%s", strings.Replace(k, "/", ".", -1))

			metric := &model.MetricValue{
				Endpoint:  hostname,
				Metric:    key,
				Value:     v,
				Timestamp: now,
				Step:      interval,
				Type:      "COUNTER",
				Tags:      "mesos,master",
			}

			fmt.Printf("%v\n", metric)
		}
	}
}
