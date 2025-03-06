package loc

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/set2105/whatTheTimeBot/internal/config"
)

type LocalTimeMap struct {
	Position map[int]*LocalTime `json:"map"`
	Name     map[string]int     `json:"-"`
	Time     time.Time          `json:"-"`

	mx      sync.Mutex
	storage *config.ConfigStorage
}

func InitLocalTimeMap(data []byte) (*LocalTimeMap, error) {
	ltm := LocalTimeMap{}
	ltm.Position = map[int]*LocalTime{}
	ltm.Name = map[string]int{}

	if data != nil {
		err := json.Unmarshal(data, &ltm)
		if err != nil {
			return &ltm, err
		}
	}

	for i, lt := range ltm.Position {
		lt.Validate()
		if !lt.Valid() {
			ltm.Delete(lt.Name)
		} else {
			lt.SetTime(&ltm.Time)
			ltm.Name[lt.Name] = i
		}
	}
	ltm.FormatPosition()
	return &ltm, nil
}

func (ltm *LocalTimeMap) GetByPosition(p string) (*LocalTime, error) {
	i, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}
	ltm.mx.Lock()
	if l, exists := ltm.Position[i]; !exists {
		ltm.mx.Unlock()
		return nil, fmt.Errorf("not found locale with position %d", i)
	} else {
		ltm.mx.Unlock()
		return l, nil
	}
}

func (ltm *LocalTimeMap) FormatPosition() {
	positionList := []int{}
	newPosition := map[int]*LocalTime{}
	newName := map[string]int{}
	ltm.mx.Lock()
	for k := range ltm.Position {
		positionList = append(positionList, k)
	}
	slices.Sort(positionList)
	for i := range positionList {
		lt := ltm.Position[positionList[i]]
		newPosition[i] = lt
		newName[lt.Name] = i
	}
	ltm.Position = newPosition
	ltm.Name = newName
	ltm.mx.Unlock()
}

func (ltm *LocalTimeMap) Values() []LocalTime {
	positionList := []int{}
	ltm.mx.Lock()
	for k := range ltm.Position {
		positionList = append(positionList, k)
	}
	slices.Sort(positionList)
	res := make([]LocalTime, len(positionList))
	for i, k := range positionList {
		res[i] = *ltm.Position[k]
	}
	ltm.mx.Unlock()
	return res
}

func (ltm *LocalTimeMap) Add(lt *LocalTime) {
	if lt == nil || !lt.Valid() {
		return
	}
	ltm.FormatPosition()
	ltm.mx.Lock()
	len := len(ltm.Position)
	ltm.Position[len] = lt
	if i, exists := ltm.Name[lt.Name]; exists {
		delete(ltm.Position, i)
	}
	ltm.Name[lt.Name] = len
	lt.SetTime(&ltm.Time)
	ltm.mx.Unlock()
}

func (ltm *LocalTimeMap) Delete(name string) {
	ltm.mx.Lock()
	if i, exists := ltm.Name[name]; exists {
		delete(ltm.Position, i)
		delete(ltm.Name, name)
	}
	ltm.mx.Unlock()
}

func (ltm *LocalTimeMap) JSON() ([]byte, error) {
	ltm.mx.Lock()
	b, err := json.Marshal(ltm)
	ltm.mx.Unlock()
	return b, err
}
