package idsyncmap

import (
	"maps"
	"sync"
)

type IDSyncMapIDTypes interface {
	int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint
}

type IDSyncMap[IDType IDSyncMapIDTypes, ItemType any] interface {
	GetAll() map[IDType]ItemType
	Add(item ItemType) (id IDType)
	Rm(id IDType)
}

type idSyncMap[IDType IDSyncMapIDTypes, ItemType any] struct {
	IDSyncMap[IDType, ItemType]
	Lock   *sync.RWMutex
	Data   map[IDType]ItemType
	NextID IDType
}

func (idsyncmap *idSyncMap[IDType, ItemType]) GetAll() map[IDType]ItemType {
	dataCopy := make(map[IDType]ItemType, len(idsyncmap.Data))
	idsyncmap.Lock.RLock()
	maps.Copy(dataCopy, idsyncmap.Data)
	idsyncmap.Lock.RUnlock()
	return dataCopy
}

func (idsyncmap *idSyncMap[IDType, ItemType]) Add(item ItemType) (id IDType) {
	idsyncmap.Lock.Lock()
	id = idsyncmap.NextID
	idsyncmap.NextID++
	idsyncmap.Data[id] = item
	idsyncmap.Lock.Unlock()
	return id
}

func (idsyncmap *idSyncMap[IDType, ItemType]) Rm(id IDType) {
	idsyncmap.Lock.Lock()
	delete(idsyncmap.Data, id)
	idsyncmap.Lock.Unlock()
}

func NewIDSyncMap[IDType IDSyncMapIDTypes, ItemType any]() IDSyncMap[IDType, ItemType] {
	return &idSyncMap[IDType, ItemType]{
		Lock:   &sync.RWMutex{},
		Data:   map[IDType]ItemType{},
		NextID: 0,
	}
}
