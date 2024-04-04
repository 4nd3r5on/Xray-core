package xray_trojan_callbacks

import (
	xray_common_callbacks "github.com/4nd3r5on/Xray-core/common/callbacks"
	"github.com/4nd3r5on/Xray-core/common/idsyncmap"
	"github.com/4nd3r5on/Xray-core/common/session"
)

type ServerCallbackManager struct {
	CbsOnProcess idsyncmap.IDSyncMap[xray_common_callbacks.OnProcess]
}

func (cm *ServerCallbackManager) ExecOnProcess(inbound *session.Inbound) (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec(inbound)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewServerCallbackManager() *ServerCallbackManager {
	return &ServerCallbackManager{
		CbsOnProcess: idsyncmap.NewIDSyncMap[xray_common_callbacks.OnProcess](),
	}
}
