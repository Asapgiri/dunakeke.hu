package dbase

import (
	"sync"
)

var muxArray = map[string]*sync.Mutex{}

func muxLock(id string) {
    mux, ok := muxArray[id]
    if !ok {
        muxArray[id] = &sync.Mutex{}
        mux = muxArray[id]
    }

    mux.Lock()
    log.Println(muxArray)
}

func muxUnlock(id string) {
    mux, ok := muxArray[id]
    if !ok {
        return
    }

    mux.Unlock()
    delete(muxArray, id)
    log.Println(muxArray)
}
