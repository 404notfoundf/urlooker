package cron

import (
	"fmt"
	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/agent/utils"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/710leo/urlooker/modules/agent/backend"
	"github.com/710leo/urlooker/modules/agent/g"
)

/*
func StartCheck() {
	t1 := time.NewTicker(time.Duration(g.Config.Web.Interval) * time.Second)
	for {
		items, err := GetItem()
		if err != nil {
			log.Println("[ERROR] ", err)
		}
		for _, item := range items {
			g.WorkerChan <- 1
			go utils.CheckTargetStatus(item)
		}
		<-t1.C
	}
}


func GetItem() ([]*dataobj.DetectedItem, error) {
	var resp dataobj.GetItemResponse
	log.Println(g.Config.IDC)
	err := backend.CallRpc("Web.GetItem", g.Config.IDC, &resp)
	if err != nil {
		return []*dataobj.DetectedItem{}, err
	}
	if resp.Message != "" {
		err := fmt.Errorf(resp.Message)
		return []*dataobj.DetectedItem{}, err
	}

	return resp.Data, err
}
*/

/*
	重新计算，需要根据相应的时间，创建定时器
*/

func Check() {
	var once sync.Once
	for {
		once.Do(func() {
			items, err := GetItemWithInterval()
			if err != nil {
				log.Println("[ERROR] ", err)
			}
			log.Println("len(items)", len(items))
			for index, item := range items {
				log.Println("index, ", index, "item, ", item)
				go func(timeDuration int) {
					t1 := time.NewTicker(time.Duration(timeDuration) * time.Second)
					for {
						select {
						case <-t1.C:
							// 获取相同时间所有的url
							data, err := GetItemWithSameInterval(timeDuration)
							if err != nil {
								log.Println("error is", err.Error())
								break
							}
							for _, i := range data {
								g.WorkerChan <- 1
								go utils.CheckTargetStatusWithInterval(i)
							}
						}
					}
				}(index)
			}
		})
	}
}

func GetItemWithInterval() (map[int][]*dataobj.DetectedItemWithInterval, error) {
	var resp dataobj.GetItemWithIntervalResponse
	err := backend.CallRpc("Web.GetItemWithInterval", g.Config.IDC, &resp)
	if err != nil {
		return map[int][]*dataobj.DetectedItemWithInterval{}, err
	}
	if resp.Message != "" {
		err := fmt.Errorf(resp.Message)
		return map[int][]*dataobj.DetectedItemWithInterval{}, err
	}
	return resp.Data, err
}

func GetItemWithSameInterval(interval int) ([]*dataobj.DetectedItemWithInterval, error) {
	// 转换为字符串数组
	var m []string
	m = append(m, g.Config.IDC)
	m = append(m, strconv.Itoa(interval))
	var resp dataobj.GetItemWithSameIntervalResponse
	err := backend.CallRpc("Web.GetItemWithSameInterval", m, &resp)
	if err != nil {
		return []*dataobj.DetectedItemWithInterval{}, err
	}
	if resp.Message != "" {
		err := fmt.Errorf(resp.Message)
		return []*dataobj.DetectedItemWithInterval{}, err
	}
	return resp.Data, err
}
