package cron

import (
	"btc_order/internal/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (j *Job) btcOrderJob() {
	for {
		ticker := time.NewTimer(2 * time.Second)
		<-ticker.C
		defer ticker.Stop()

		var btcOrders []types.BtcOrder
		err := j.db.Model(&btcOrders).
			Where("status = 'pending'").
			Order("ex_time ASC").Select()
		if err != nil {
			return
		}
		for _, btcOrder := range btcOrders {
			if btcOrder.TxId == "" && btcOrder.ExTime <= time.Now().UnixMilli() {
				btcOrder.Status = "unpaid"
				_, err = j.db.Model(&btcOrder).WherePK().Update()
				if err != nil {
					return
				}
			} else if btcOrder.TxId != "" {
				url := fmt.Sprintf("https://mempool.space/api/tx/%s", btcOrder.TxId)
				resp, err := http.Get(url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				var bodys map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&bodys)
				if err != nil {
					return
				}
				status := bodys["status"].(map[string]interface{})
				c := status["confirmed"].(bool)
				if c {
					btcOrder.Status = "inscribed"
					_, err = j.db.Model(&btcOrder).WherePK().Update()
					if err != nil {
						return
					}
				}
				time.Sleep(1 * time.Second)
			}

		}
	}
}
