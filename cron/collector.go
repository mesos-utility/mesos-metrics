package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/open-falcon/common/model"
	"github.com/soarpenguin/mesos-metrics/g"
)

func Collect() {
	if !g.Config().Transfer.Enable {
		return
	}

	if g.Config().Transfer.Addr == "" {
		return
	}

	go collect(g.Config().Services)
}

func collect(sev []*g.ServiceConfig) {

	// start collect data for mesos cluster.
	for {
	REST:
		var interval int64 = g.Config().Transfer.Interval
		time.Sleep(time.Duration(interval) * time.Second)
		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		mvs := []*model.MetricValue{}
		for _, srv := range g.Config().Services {
			if !srv.Enable {
				if g.Config().Debug {
					log.Printf("[Notice]: %s not enabled!!!", srv.Type)
				}
				continue
			}

			addr := srv.Apiurl
			srvtype := srv.Type
			resp, err := http.Get(addr)
			if err != nil {
				log.Println("get mesos metric data fail", err)
				continue
			}
			defer resp.Body.Close()

			// read json http response
			jsonData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("read mesos metric fail", err)
				continue
			}

			var f interface{}
			err = json.Unmarshal(jsonData, &f)
			if err != nil {
				log.Println("Unmarshal metric data fail", err)
				continue
			}

			now := time.Now().Unix()
			m := f.(map[string]interface{})
			for k, v := range m {
				key := fmt.Sprintf("mesos.%s", strings.Replace(k, "/", ".", -1))

				metric := &model.MetricValue{
					Endpoint:  hostname,
					Metric:    key,
					Value:     v,
					Timestamp: now,
					Step:      interval,
					Type:      "GAUGE",
					Tags:      srvtype,
				}

				mvs = append(mvs, metric)
				//fmt.Printf("%v\n", metric)
			}
		}
		g.SendToTransfer(mvs)
	}
}
