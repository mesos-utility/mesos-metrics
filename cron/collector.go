package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	//for _, v := range funcs.Mappers {
	//	go collect(int64(v.Interval), v.Fs)
	//}
	go collect()
}

//func collect(sec int64, fns []func() []*model.MetricValue) {
func collect() {

	gcfg := g.Config()
	addr := gcfg.Master.Apiurl
	hostname, _ := g.Hostname()

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
