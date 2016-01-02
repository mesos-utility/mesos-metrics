package funcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mesos-utility/mesos-metrics/g"
	"github.com/open-falcon/common/model"
)

func CollectMetrics(hostname string, srvCfg *g.ServiceConfig) []*model.MetricValue {

	var interval int64 = g.Config().Transfer.Interval
	var stats = map[string]string{}

	addr := srvCfg.Apiurl
	srvtype := srvCfg.Type
	resp, err := http.Get(addr)

	if srvtype == "master" {
		stats = statsMaster
	} else if srvtype == "slave" {
		stats = statsSlave
	} else {
		log.Println("not support server type!!!", err)
		return nil
	}

	if err != nil {
		log.Println("get mesos metric data fail", err)
		return nil
	}
	defer resp.Body.Close()

	// read json http response
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read mesos metric fail", err)
		return nil
	}

	var f interface{}
	err = json.Unmarshal(jsonData, &f)
	if err != nil {
		log.Println("Unmarshal metric data fail", err)
		return nil
	}

	now := time.Now().Unix()
	m := f.(map[string]interface{})
	mvs := []*model.MetricValue{}
	for k, v := range m {
		mtype := "GAUGE"
		if t, ok := stats[k]; ok {
			mtype = t
		}

		key := fmt.Sprintf("mesos.%s", strings.Replace(k, "/", ".", -1))

		metric := &model.MetricValue{
			Endpoint:  hostname,
			Metric:    key,
			Value:     v,
			Timestamp: now,
			Step:      interval,
			Type:      mtype,
			Tags:      srvtype,
		}

		mvs = append(mvs, metric)
	}

	return mvs
}
