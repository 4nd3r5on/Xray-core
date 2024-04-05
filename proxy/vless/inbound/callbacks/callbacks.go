package xray_vless_inbound_callbacks

import (
	xray_common_callbacks "github.com/4nd3r5on/Xray-core/common/callbacks"
	"github.com/4nd3r5on/Xray-core/common/idsyncmap"
	"github.com/4nd3r5on/Xray-core/common/session"
	"github.com/4nd3r5on/Xray-core/features/policy"
)

type OnProcessStart interface {
	Exec(sessionPolicy *policy.Session) error
}

type CallbackManager struct {
	CbsOnProcess      idsyncmap.IDSyncMap[xray_common_callbacks.OnProcess]
	CbsOnProcessStart idsyncmap.IDSyncMap[OnProcessStart]
}

func (cm *CallbackManager) ExecOnProcess(inbound *session.Inbound) (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec(inbound)
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
		CbsOnProcess: idsyncmap.NewIDSyncMap[xray_common_callbacks.OnProcess](),
	}
}
