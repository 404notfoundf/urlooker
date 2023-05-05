package cron

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/710leo/urlooker/dataobj"

	"github.com/710leo/urlooker/modules/agent/backend"
	"github.com/710leo/urlooker/modules/agent/g"
	"github.com/710leo/urlooker/modules/agent/utils"
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
*/

func Check() {
	var once sync.Once
	for {
		once.Do(func() {
			items, err := GetItemInterval()
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
							data, err := GetItemSameInterval(timeDuration)
							if err != nil {
								log.Println("error is", err.Error())
								break
							}
							for _, i := range data {
								g.WorkerChan <- 1
								go utils.CheckTargetStatus(i)
							}
						}
					}
				}(index)
			}
		})
	}
}


func GetItemSameInterval(interval int) (data []*dataobj.DetectedItemWithInterval, err error) {
	itemInterval, err := GetItemInterval()
	if err != nil {
		return []*dataobj.DetectedItemWithInterval{}, err
	}
	if _, exists := itemInterval[interval]; exists {
		return itemInterval[interval], nil
	} else {
		return []*dataobj.DetectedItemWithInterval{}, errors.New("not found")
	}
}

func GetItemInterval() (map[int][]*dataobj.DetectedItemWithInterval, error) {
	// 将得到的数组分成合适的份数
	items, err := GetItem()
	if err != nil {
		return map[int][]*dataobj.DetectedItemWithInterval{}, err
	}
	m := make(map[int][]*dataobj.DetectedItemWithInterval)
	for _, item := range items {
		decetcedItem := newDetectedItem(item)
		m[decetcedItem.Interval] = append(m[decetcedItem.Interval], decetcedItem)
	}
	return m, nil
}


func newDetectedItem (item *dataobj.DetectedItem) *dataobj.DetectedItemWithInterval{
	log.Println("item", item)
	var force bool
	var interval int
	if len(g.Config.UrlInterval) != 0 {
		for _, u := range g.Config.UrlInterval {
			var isOK bool
			for _, index := range u.Url {
				if index == item.Target {
					isOK = true
					break
				}
			}
			if isOK {
				interval = u.Interval
				force = true
				break
			}
		}
	}
	var n *dataobj.DetectedItemWithInterval
	n.Target = item.Target
	n.Tag = item.Tag
	n.Idc = item.Idc
	n.Creator = item.Creator
	n.Sid = item.Sid
	n.Keywords = item.Keywords
	n.Data = item.Data
	n.Endpoint = item.Endpoint
	n.Timeout = item.Timeout
	n.Header = item.Header
	n.PostData = item.PostData
	n.Method = item.Method
	n.Domain = item.Domain
	n.ExpectCode = item.ExpectCode
	if force {
		n.Interval = interval
	} else {
		n.Interval = g.Config.Web.Interval
	}
	return n
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
