package vmess_inbound_callbacks

import (
	callbacksCommon "github.com/xtls/xray-core/common/callbacks"
	"github.com/xtls/xray-core/common/idsyncmap"
	"github.com/xtls/xray-core/common/session"
	"github.com/xtls/xray-core/features/policy"
)

type OnProcessStart func(sessionPolicy *policy.Session) error

type CallbackManager struct {
	CbsOnProcess      idsyncmap.IDSyncMap[int, callbacksCommon.OnProcess]
	CbsOnProcessStart idsyncmap.IDSyncMap[int, OnProcessStart]
}

func (cm *CallbackManager) ExecOnProcess(inbound *session.Inbound) (id int, err error) {
	for id, callback := range cm.CbsOnProcess.GetAll() {
		err = callback(inbound)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func (cm *CallbackManager) ExecOnProcessStart(sessionPolicy *policy.Session) (id int, err error) {
	for id, callback := range cm.CbsOnProcessStart.GetAll() {
		err = callback(sessionPolicy)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewCallbackManager() *CallbackManager {
	return &CallbackManager{
		CbsOnProcess:      idsyncmap.NewIDSyncMap[int, callbacksCommon.OnProcess](),
		CbsOnProcessStart: idsyncmap.NewIDSyncMap[int, OnProcessStart](),
	}
}
