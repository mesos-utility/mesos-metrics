package cron

import (
	"fmt"
	"log"
	"time"

	"github.com/mesos-utility/mesos-metrics/funcs"
	"github.com/mesos-utility/mesos-metrics/g"
	"github.com/open-falcon/common/model"
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

func collect(SrvCfgs []*g.ServiceConfig) {

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
		for _, srv := range SrvCfgs {
			if !srv.Enable {
				if g.Config().Debug {
					log.Printf("[Notice]: %s not enabled!!!", srv.Type)
				}
				continue
			}

			mvs = funcs.CollectMetrics(hostname, srv)
		}

		g.SendToTransfer(mvs)
		//fmt.Printf("%v\n", mvs)
	}
}
