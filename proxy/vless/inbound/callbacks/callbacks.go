package xray_vless_inbound_callbacks

import (
	"github.com/4nd3r5on/Xray-core/common/idsyncmap"
	"github.com/4nd3r5on/Xray-core/features/policy"
	"github.com/4nd3r5on/Xray-core/proxy/vless"
)

type (
	OnProcess struct {
		Exec func(account *vless.MemoryAccount) error
	}
	OnProcessStart struct {
		Exec func(sessionPolicy *policy.Session) error
	}
)

type CallbackManager struct {
	CbsOnProcess      idsyncmap.IDSyncMap[OnProcess]
	CbsOnProcessStart idsyncmap.IDSyncMap[OnProcessStart]
}

func (cm *CallbackManager) ExecOnProcess(account *vless.MemoryAccount) (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec(account)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func (cm *CallbackManager) ExecOnProcessStart(sessionPolicy *policy.Session) (id int32, err error) {
	for id, callback := range cm.CbsOnProcessStart.Get() {
		err = callback.Exec(sessionPolicy)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewCallbackManager() *CallbackManager {
	return &CallbackManager{
		CbsOnProcess: idsyncmap.NewIDSyncMap[OnProcess](),
	}
}
