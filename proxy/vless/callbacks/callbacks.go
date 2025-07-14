package vless_callbacks

import (
	"github.com/xtls/xray-core/common/callbacks"
	callbacksCommon "github.com/xtls/xray-core/common/callbacks"
	"github.com/xtls/xray-core/common/idsyncmap"
	"github.com/xtls/xray-core/common/session"
	"github.com/xtls/xray-core/features/policy"
)

type OnProcessStart func(sessionPolicy *policy.Session) error

type InboundCallbackManager struct {
	CbsOnProcess      idsyncmap.IDSyncMap[int, callbacks.InboundOnProcess]
	CbsOnProcessStart idsyncmap.IDSyncMap[int, OnProcessStart]
}

func (cm *InboundCallbackManager) ExecOnProcess(inbound *session.Inbound) (id int, err error) {
	for id, callback := range cm.CbsOnProcess.GetAll() {
		err = callback(inbound)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func (cm *InboundCallbackManager) ExecOnProcessStart(sessionPolicy *policy.Session) (id int, err error) {
	for id, callback := range cm.CbsOnProcessStart.GetAll() {
		err = callback(sessionPolicy)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewInboundCallbackManager() *InboundCallbackManager {
	return &InboundCallbackManager{
		CbsOnProcess:      idsyncmap.NewIDSyncMap[int, callbacksCommon.InboundOnProcess](),
		CbsOnProcessStart: idsyncmap.NewIDSyncMap[int, OnProcessStart](),
	}
}
