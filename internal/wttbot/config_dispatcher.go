package wttbot

import (
	"fmt"
	"sync"
	"time"

	"github.com/set2105/whatTheTimeBot/internal/config"
	"github.com/set2105/whatTheTimeBot/internal/loc"
)

type ConfigDispatcher struct {
	storage *config.ConfigStorage

	configMap  map[int]*loc.LocalTimeMap
	lastMsgMap map[int]time.Time

	mx sync.Mutex
}

func InitConfigDispatcher(storage *config.ConfigStorage) *ConfigDispatcher {
	cd := ConfigDispatcher{}
	cd.configMap = map[int]*loc.LocalTimeMap{}
	cd.lastMsgMap = map[int]time.Time{}
	cd.storage = storage
	return &cd
}

func (cd *ConfigDispatcher) StartClear() {
	for true {
		time.Sleep(time.Hour)
		cd.Clear()
	}
}

func (cd *ConfigDispatcher) Clear() {
	checkTime := time.Now().Add(-time.Hour)
	cd.mx.Lock()
	for i, lastMsgTime := range cd.lastMsgMap {
		if lastMsgTime.After(checkTime) {
			delete(cd.configMap, i)
			delete(cd.lastMsgMap, i)
		}
	}
	cd.mx.Unlock()
}

func (cd *ConfigDispatcher) Get(chatId int) (ltm *loc.LocalTimeMap, err error) {
	if conf, exists := cd.configMap[chatId]; exists {
		cd.mx.Lock()
		cd.lastMsgMap[chatId] = time.Now()
		cd.mx.Unlock()
		return conf, nil
	}
	b, err := cd.storage.Read(fmt.Sprint(chatId))
	if err != nil {
		ltm, err = loc.InitLocalTimeMap(nil)
		if err != nil {
			return nil, err
		}
	} else {
		ltm, err = loc.InitLocalTimeMap(b)
		if err != nil {
			return nil, err
		}
	}
	cd.mx.Lock()
	cd.configMap[chatId] = ltm
	cd.lastMsgMap[chatId] = time.Now()
	cd.mx.Unlock()
	return ltm, nil
}

func (cd *ConfigDispatcher) Save(chatId int) error {
	b, err := cd.configMap[chatId].JSON()
	if err != nil {
		return err
	}
	err = cd.storage.Save(fmt.Sprint(chatId), b)
	if err != nil {
		return err
	}
	return nil
}
