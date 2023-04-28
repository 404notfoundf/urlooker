package g

import (
	"sync"

	"github.com/710leo/urlooker/dataobj"
)

type DetectedItemSafeMap struct {
	sync.RWMutex
	M map[string][]*dataobj.DetectedItem
}

type DetectedItemWithIntervalSafeMap struct {
	sync.RWMutex
	M map[string]map[int][]*dataobj.DetectedItemWithInterval
}

var (
	DetectedItemMap             = &DetectedItemSafeMap{M: make(map[string][]*dataobj.DetectedItem)}
	DetectedItemWithIntervalMap = &DetectedItemWithIntervalSafeMap{M: make(map[string]map[int][]*dataobj.DetectedItemWithInterval)}
)

func (this *DetectedItemSafeMap) Get(key string) ([]*dataobj.DetectedItem, bool) {
	this.RLock()
	defer this.RUnlock()
	ipItem, exists := this.M[key]
	return ipItem, exists
}

func (this *DetectedItemSafeMap) GetAll() map[string][]*dataobj.DetectedItem {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *DetectedItemSafeMap) Set(detectedItemMap map[string][]*dataobj.DetectedItem) {
	this.Lock()
	defer this.Unlock()
	this.M = detectedItemMap
}

func (that *DetectedItemWithIntervalSafeMap) Get(key string) (map[int][]*dataobj.DetectedItemWithInterval, bool) {
	that.RLock()
	defer that.RUnlock()
	ipItem, exists := that.M[key]
	return ipItem, exists
}

func (that *DetectedItemWithIntervalSafeMap) GetAll() map[string]map[int][]*dataobj.DetectedItemWithInterval {
	that.RLock()
	defer that.RUnlock()
	return that.M
}

func (that *DetectedItemWithIntervalSafeMap) Set(detectedItemMap map[string]map[int][]*dataobj.DetectedItemWithInterval) {
	that.Lock()
	defer that.Unlock()
	that.M = detectedItemMap
}
