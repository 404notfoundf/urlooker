package cron

import (
	g2 "github.com/710leo/urlooker/modules/agent/g"
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
		// err := getDetectedItem()
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
	detectedItemMap := make(map[string]map[int][]*dataobj.DetectedItem)
	stras, err := model.GetAllStrategyByCron()
	if err != nil {
		log.Println("get strategies error:", err)
	}
	for _, s := range stras {
		detectedItem := newDetectedItem(s)
		idc := detectedItem.Idc
		d := dataobj.DetectedItem{
			Idc:        idc,
			Target:     detectedItem.Target,
			Creator:    detectedItem.Creator,
			Sid:        detectedItem.Sid,
			Keywords:   detectedItem.Keywords,
			Data:       detectedItem.Data,
			Tag:        detectedItem.Tag,
			Endpoint:   detectedItem.Endpoint,
			ExpectCode: detectedItem.ExpectCode,
			Timeout:    detectedItem.Timeout,
			Header:     detectedItem.Header,
			PostData:   detectedItem.PostData,
			Method:     detectedItem.Method,
			Domain:     detectedItem.Domain,
		}
		if _, exists := detectedItemMap[idc][detectedItem.Interval]; exists {
			// 修改类型
			detectedItemMap[idc][detectedItem.Interval] = append(detectedItemMap[idc][detectedItem.Interval], &d)
		} else {
			detectedItemMap[idc][detectedItem.Interval] = []*dataobj.DetectedItem{&d}
		}
	}
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
	for _, u := range g.Config.UrlInterval {
		if u.Url == s.Url {
			force = true
			interval = u.Interval
			break
		}
	}
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
	if force {
		detectedItem.Interval = interval
	} else {
		detectedItem.Interval = g2.Config.Web.Interval
	}
	return detectedItem
}
