package trojan

import (
	"github.com/4nd3r5on/Xray-core/common/idsyncmap"
	"github.com/4nd3r5on/Xray-core/common/session"
	"github.com/4nd3r5on/Xray-core/features/policy"
)

type (
	OnProcess struct {
		Exec func(inbound *session.Inbound, sessionPolicy policy.Session) error
	}
)

type CallbackManager struct {
	CbsOnProcess idsyncmap.IDSyncMap[OnProcess]
}

func (cm *CallbackManager) ExecOnProcess(inbound *session.Inbound, sessionPolicy policy.Session) (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec(inbound, sessionPolicy)
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
