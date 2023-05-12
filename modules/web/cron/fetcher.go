package cron

import (
	"log"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/model"
	"github.com/710leo/urlooker/modules/web/utils"
)

func GetDetectedItem() {
	t1 := time.NewTicker(time.Duration(60) * time.Second)
	for {
		err := getDetectedItemWihInterval()
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		<-t1.C
	}
}

/*
func getDetectedItem() error {
	detectedItemMap := make(map[string][]*dataobj.DetectedItem)
	stras, err := model.GetAllStrategyByCron()
	if err != nil {
		log.Println("get strategies error:", err)
		return err
	}

	for _, s := range stras {
		detectedItem := newDetectedItem(s)
		idc := detectedItem.Idc
		if _, exists := detectedItemMap[idc]; exists {
			detectedItemMap[idc] = append(detectedItemMap[idc], &detectedItem)
		} else {
			detectedItemMap[idc] = []*dataobj.DetectedItem{&detectedItem}
		}
	}
	g.DetectedItemMap.Set(detectedItemMap)
	return nil
}
*/

func getDetectedItemWihInterval() error {
	detectedItemMap := make(map[string]map[int][]*dataobj.DetectedItemWithInterval)
	strateges, err := model.GetAllStrategyByCron()
	if err != nil {
		log.Println("get strategies error:", err)
		return err
	}
	for _, s := range strateges {
		log.Println("!!!!star", s)
		detectedItem := newDetectedItem(s)
		idc := detectedItem.Idc
		if _, exists := detectedItemMap[idc][detectedItem.Interval]; exists {
			// 修改类型
			detectedItemMap[idc][detectedItem.Interval] = append(detectedItemMap[idc][detectedItem.Interval], &detectedItem)
		} else {
			if detectedItemMap[idc] == nil {
				// 注意在二维的map中，一维的map也需要使用make来构造一下
				detectedItemMap[idc] = make(map[int][]*dataobj.DetectedItemWithInterval)
			}
			detectedItemMap[idc][detectedItem.Interval] = []*dataobj.DetectedItemWithInterval{&detectedItem}
		}
	}

	g.DetectedItemWithIntervalMap.Set(detectedItemMap)
	return nil
}

func newDetectedItem(s *model.Strategy) dataobj.DetectedItemWithInterval {
	_, domain, _, _ := utils.ParseUrl(s.Url)
	idc := s.Idc
	if idc == "" {
		idc = g.Config.IDC[0]
	}
	var force bool
	var interval int
	// 判断是否为空
	log.Println("len g.urlinterval is", len(g.Config.UrlInterval))
	// log.Println("!!!g.config.urlinterval", g.Config.UrlInterval)
	// 比较不同情况下的url是否相同

	if len(g.Config.UrlInterval) != 0 {
		for _, u := range g.Config.UrlInterval {
			var isOK bool
			for _, index := range u.Url {
				if index == s.Url {
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

	/*
		if len(g.Config.UrlInterval) != 0 {
			for _, u := range g.Config.UrlInterval {
				if u.Url == s.Url {
					force = true
					interval = u.Interval
					break
				}
			}
		}
	*/
	var detectedItem dataobj.DetectedItemWithInterval
	detectedItem.Idc = idc
	detectedItem.Target = s.Url
	detectedItem.Creator = s.Creator
	detectedItem.Sid = s.Id
	detectedItem.Keywords = s.Keywords
	detectedItem.Data = s.Data
	detectedItem.Tag = s.Tag
	detectedItem.Endpoint = s.Endpoint
	detectedItem.Timeout = s.Timeout
	detectedItem.Header = s.Header
	detectedItem.PostData = s.PostData
	detectedItem.Method = s.Method
	detectedItem.Domain = domain
	detectedItem.ExpectCode = s.ExpectCode
	if force {
		detectedItem.Interval = interval
	} else {
		// TODO： 暂时写死一下
		detectedItem.Interval = 60
		//detectedItem.Interval = g2.Config.Web.Interval
	}
	log.Println("detecterItem is", detectedItem)
	return detectedItem
}
