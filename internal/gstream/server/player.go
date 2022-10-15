package server

import "sync"

//type PlayerAll struct {
//	Data       map[string]Data
//	PlayerLock *sync.RWMutex
//}

var PlayerAll map[string]*Data

var PLock *sync.RWMutex

type Data struct {
	SessionId int32
	DeskId    int32
}

func init() {
	PlayerAll = make(map[string]*Data)
	PLock = new(sync.RWMutex)
}
