package g

import (
	"log"

	"github.com/open-falcon/common/model"
)

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
